# github.com/revdfdev/mongodb-helpers

## Usage

### Initialize your module

```
go mod init example.com/my-awesome-project
```

### Get the go

Note that you need to include the **v** in the version tag.

```
go get github.com/revdfdev/mongodb-helpers@v0.1.2
```

```go
package main

import (
    "fmt"

    "github.com/revdfdev/mongodb-helpers/database"
)

func main() {
    uri := os.GetEnv("MONGODB_URI")
    db, err := database.NewDatabaseConnection(uri)

 if err != nil {
  panic(err)
 }
}
```

### Commonly used function

```go
// Find one
singleResult, err := database.FindOne("database_name", "collection_name", bson.M{
    //filter here.
}, bson.M{})


// Find many
cursor, err := database.FindMany("database_name", "collection_name", bson.M{}, bson.M{})

// InsertOne

dto := MyAwesomeDto()

dto.Bind(&myAwesomeRequest)

insertResult, err := database.InsertOne("database_name", "collection_name", dto)


// Similarly you can use InsertMany for slices


```

### For filters if you are query with _id field make sure you do this.
```go

 objectId, _ := primitive.ObjectIdFromHex(id);
 
 singleResult, err := database.FindOne("database_name", "collection_name", bson.M{
    "_id": objectId,
}, bson.M{})
```

### custom aggregation.
```go
 o1 := bson.M{
    // your aggregation stage here
 }

 o2 := bson.M{
   // your 2nd aggregation stage here  
 }

 aggResult, err := database.CustomAggregate("database_name", "collection_name", []bson.M{o1, o2})
```

## Tagging

```
git tag v0.1.0
git push origin --tags
```
