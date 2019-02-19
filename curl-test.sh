printf "Example of posting content\n"
curl --header "Content-Type: application/json" --request POST --data '{"Content":"{\"Title\":\"Wonderful\"}"}' http://localhost:8080/api/write

printf "\n\nExample of requesting latest N posts\n"
curl http://localhost:8080/api/list/


printf "\n\nExample of a read request -- replace ABCDEF with a PostId from the latest N posts\n"
printf "curl http://localhost:8080/api/read/ABCDEF\n\n"
curl http://localhost:8080/api/read/ABCDEF

