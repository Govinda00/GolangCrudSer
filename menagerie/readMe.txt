1) extract the zip file 

Need to install goland and mysql in the local system

go run main.go

2) change the credential as per your system
file path: db/db.go
change: 
user:=root
password=Govinda@123
host=localhost:3306
db=servicaGolang

DB, err = sql.Open("mysql", "root:Govinda@123@tcp(localhost:3306)/servicaGolang")

3) your check the serve is not in use, if the change it 8XXX
file path: main.go

if err := http.ListenAndServe(":8002", r)

4)collection is attchec in the file and import the collection in your postman/any similar tool

RestAPIGoServika.postman_collection.json

5) now the run the main.go file
go run main.go

now you can perform all the operation as per collection.