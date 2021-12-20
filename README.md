# GoExercise
![alt text](golang_icon_no_bg.png)

## General Info:
This program does a distributed grep using the MapReduce paradigm, 
typical of big data problems. The **GoGrep** program returns the lines of a large text file 
given in input that match a specific regex specified (i.e., a regular expression). The 
input file is chosen by the end user.
The program is written in **Go** and uses **RPC** for the communication between client, master and
worker peers.

## Running:
To run the program, is necessary to first set up the master and the worker peers and initialize
each corresponding module.
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

### Run the master:
Open a new terminal and input the following command:
````
go run master/master.go
````

### Run the worker:
Open a new terminal and input the following command:
````
go run worker/worker.go
````
Once the master and the worker peer are up, the client also can run connecting to the master
on the random port in which it's listening for incoming connections.

### Run the client:
Open a new terminal and input the following command:
````
go run client/client.go
````

It' possible to choose a specific file to grep, and in that case the file must be present
in the directory **client/files**. Then the user must specify the regex. 
It's possible to choose multiple regexes to use for the grep on the file chosen. If any of the 
lines of the text contains the regex specified, then the master returns the set of those lines.

## MapReduce
Once the file is chosen, the master splits its content in multiple chunks. Each chunk is then
distributed to a different worker. The master spawns N workers and the number N of workers is
determined equally distributing the total lines of the files, such that workers processes
a maximum of 10 lines each.
Each worker performs the map phase on the chunk of the file received by the master, then
returns the result of the mapping to the master, that proceeds performing the following operations:
* Shuffle
* Sort
* Reduce

The master merges the results obtained by the workers and returns the subset of lines to the client.
