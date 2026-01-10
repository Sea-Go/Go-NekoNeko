#include<bits/stdc++.h>
using i64 =long long;

struct node{
    int l,r;
    int val;
};
std::unordered_map<int,node> list;
std::unordered_map<int,int> val1; 
int head=0;
int maxsize;
int get(int idx){
    if(val1.count(idx)==0){
        return 0;
    }
    int id=val1[idx];
    if(list[id].l==0){
        return list[id].val;
    }
    list[list[id].l].r=list[id].r;
    list[list[id].r].l=list[id].l;
    list[id].r=head;
    list[id].l=list[head].l;
    list[head].l=id;
    list[list[id].l].r=id;
    head=id;
}
int cur=0;
void put(int idx,int val){
    if(cur==maxsize){
        int wei=list[head].l;
        list[wei].val=val;
        head=wei;
        val1[idx]=wei;
    }
    else if(cur==0){
        list[cur].val=val;
        list[cur].r=-1;
        list[cur].l=-1;
        val1[idx]=cur;;
        cur++;

    }
    else if(cur==1){
        list[head].l=cur;
        list[head].r=cur;
        list[cur].l=head;
        list[cur].r=head;
        list[cur].val=val;
        head=cur;
        val1[idx]=cur;;
        cur++;
    }
    else{
        list[cur].r=head;
        list[cur].l=list[head].l;
        list[cur].val=val;
        int nexthead=cur;
        list[head].l=cur;
        list[list[cur].l].r=cur;
        head=nexthead;
        val1[idx]=cur;
        cur++;
    }
    
}
void print(){
    std::cout<<"print"<<'\n';
    int c=head;
    int w=0;
    while(list[c].r!=-1&&list[c].r!=head){
        std::cout<<c<<" "<<list[c].val <<" "<<list[c].l<<" "<<list[c].r<<std::endl;
        w++;
        c=list[c].r;
    }
    std::cout<<c<<" "<<list[c].val<<std::endl;
    return;

}
int main(){
    std::ios::sync_with_stdio(0);
    std::cin.tie(0);
    maxsize=3;
    std::string s;
    while(1){
        std::cin>>s;
        if(s=="get"){
            int idx;
            std::cin>>idx;
            std::cout<<get(idx)<<'\n';
        }
        else if(s=="put"){
            int idx,val;
            std::cin>>idx>>val;
            put(idx,val);
        }
        else{
            break;
        }
        print();
    }

}