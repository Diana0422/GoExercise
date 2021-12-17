# GoExercise
![alt text](golang_icon_no_bg.png)

##Running:
If the file **go.mod** is not in the repository **GoExercise**, then:
* Set the modules of the repository:
````
go mod init GoExercise
````
* Add module requirements and sums:
````
go mod tidy
````
Now can run the code.

### Run the client:
````
go run client/client.go
````

### Run the master:
````
go run server/master/main.go
````

