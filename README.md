# GolangCrudSer
used golang for simple crud and used validation

clone: 
https://github.com/Govinda00/GolangCrudSer.git

1) extract the zip file
Need to install goland and mysql in the local system
go run main.go

3) change the credential as per your system
file path: db/db.go
change:
user:=root
password=Govinda@123
host-localhost:3306
db=servicaGolang
DB, err = sql.Open('mysq",
"root:Govinda@123@tcp(localhost:3306)/servicaGolang")


5) your check the serve is not in use, if the change it 8XXX
file path: main.go
if err := http.ListenAndServe(":8002", r)


4)collection is attached in the file and import the collection
in your postman/any similar tool
RestAPIGoServika.postman_collection1.json


7) now the run the main.go file
* go run main.go
now you can perform all the operation as per collection.
