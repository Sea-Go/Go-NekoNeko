
### Test Script for Favorite System

# 1. Create Folder
curl -X POST http://localhost:8890/api/v1/favorite/folder/create \
  -H "Content-Type: application/json" \
  -d '{"name": "Tech Blog", "is_public": true}'

# 2. List Folders (Assuming user_id=0 or extracted from token, but wait, we need auth...)
# Since Auth is enabled (jwt: Auth), we cannot easily curl without a valid token.
# I will temporarily disable Auth in .api file to run this test script, then revert it.
