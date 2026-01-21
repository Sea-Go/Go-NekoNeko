package snowflake

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
    sf    *snowflake.Node
    once  sync.Once
    sfErr error
)

// Init 初始化雪花ID生成器
func Init() {
    once.Do(func() {
        nodeId, err := generateNodeId()
        if err != nil {
            sfErr = fmt.Errorf("generate node id failed: %w", err)
            return
        }

        // Snowflake nodeId 必须在 [0, 1023] 范围内（10位）
        nodeId = nodeId % 1024

        sf, sfErr = snowflake.NewNode(int64(nodeId))
        if sfErr == nil {
            logx.Infof("Snowflake initialized with nodeId: %d", nodeId)
        }
    })
}

// GetID 生成全局唯一ID
func GetID() (int64, error) {
    if sfErr != nil {
        return 0, sfErr
    }
    return sf.Generate().Int64(), nil
}

// generateNodeId 生成唯一 nodeId
func generateNodeId() (uint16, error) {
    // 优先从环境变量指定（用于调试或固定分配）
    if idStr := os.Getenv("SNOWFLAKE_NODE_ID"); idStr != "" {
        if id, err := strconv.Atoi(idStr); err == nil && id >= 0 {
            return uint16(id), nil
        }
    }

    // dev方案：服务名 + 进程ID（PID）实现单机下唯一
    serviceName := os.Getenv("SERVICE_NAME")
    if serviceName == "" {
        serviceName = "unknown-service"
    }
    pid := os.Getpid()
    data := fmt.Sprintf("%s-%d", serviceName, pid)

    /*
    // 上线方案：IP + 端口
    // 问题：在单机多实例测试时，所有服务可能共享相同 IP（如 127.0.0.1），
    //      若端口未显式设置或冲突，仍可能导致 nodeId 重复。
    // 使用场景：仅适用于每服务独占 IP 或严格端口隔离的生产环境。
    
    ip := getOutboundIP().String()
    port := os.Getenv("PORT")
    if port == "" {
        port = "0"
    }
    data := ip + ":" + port
    */

    hash := sha256.Sum256([]byte(data))
    return binary.BigEndian.Uint16(hash[:2]), nil
}

// getOutboundIP 获取出口IP（辅助函数，当前未使用）
// func getOutboundIP() net.IP {
//     conn, err := net.Dial("udp", "8.8.8.8:80")
//     if err != nil {
//         addrs, _ := net.InterfaceAddrs()
//         for _, addr := range addrs {
//             if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
//                 return ipnet.IP
//             }
//         }
//         return net.IPv4(127, 0, 0, 1)
//     }
//     defer conn.Close()
//     localAddr := conn.LocalAddr().(*net.UDPAddr)
//     return localAddr.IP
// }