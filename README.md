## create user
curl -d "user_id=100" http://localhost:3000/users

## get information about all users
curl http://localhost:3000/users

## create chat room
curl -d "user1=100&user2=200" http://localhost:3000/chatrooms

## get information about chat rooms including specific user
curl http://localhost:3000/chatrooms/1
