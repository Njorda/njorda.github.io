package main

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/njorda.github.io/go/CloudSQL/parser"
	"google.golang.org/api/iterator"

	"github.com/chzyer/readline"
)

//var columns = []string{"name", "size", "timeCreated", "timeUpdated", "storageClass", "owner", "contentType", "contentEncoding", "contentDisposition", "retentionTime", "updated"}

func handleInput(ctx context.Context, client *storage.Client, input string) error {
	query, err := parser.NewParser(input).ParseQuery()
	if err != nil {
		return err
	}

	// Name is the only value we have ...
	objects, err := ListObjects(ctx, client, query.From, query.Where.Value, query.Select)
	if err != nil {
		return err
	}
	for _, obj := range objects {
		fmt.Println(obj)
	}
	return nil
}

// CreateClient initializes a new Google Cloud Storage client
func CreateClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetObjects lists objects in a given bucket, optionally filtered by a prefix
func GetObjects(ctx context.Context, client *storage.Client, bucketName, objet string, selected []string) ([]string, error) {
	attrs, err := client.Bucket(bucketName).Object(objet).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	var objects []string
	out := map[string]string{}
	out["name"] = attrs.Name
	out["size"] = fmt.Sprintf("%d", attrs.Size)
	out["timeCreated"] = attrs.Created.String()
	out["timeUpdated"] = attrs.Updated.String()
	out["storageClass"] = attrs.StorageClass
	out["owner"] = attrs.Owner
	out["contentType"] = attrs.ContentType
	out["contentEncoding"] = attrs.ContentEncoding
	out["contentDisposition"] = attrs.ContentDisposition
	out["retentionTime"] = attrs.RetentionExpirationTime.GoString()
	out["updated"] = attrs.Updated.String()
	for _, column := range selected {
		objects = append(objects, out[column])
	}
	return objects, nil
}

// ListObjects lists objects in a given bucket, optionally filtered by a prefix
func ListObjects(ctx context.Context, client *storage.Client, bucketName, prefix string, selected []string) ([]string, error) {
	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{Prefix: prefix})
	var objects []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		// Here we could do something with reflect to get all the stuff out!
		// This would be the way.
		out := map[string]string{}
		out["name"] = attrs.Name
		out["size"] = fmt.Sprintf("%d", attrs.Size)
		out["timeCreated"] = attrs.Created.String()
		out["timeUpdated"] = attrs.Updated.String()
		out["storageClass"] = attrs.StorageClass
		out["owner"] = attrs.Owner
		out["contentType"] = attrs.ContentType
		out["contentEncoding"] = attrs.ContentEncoding
		out["contentDisposition"] = attrs.ContentDisposition
		out["retentionTime"] = attrs.RetentionExpirationTime.GoString()
		out["updated"] = attrs.Updated.String()
		for _, column := range selected {
			objects = append(objects, out[column])
		}
	}
	return objects, nil
}

func main() {
	ctx := context.Background()
	client, err := CreateClient(ctx) // Assuming CreateClient is a function you've defined
	if err != nil {
		panic(err)
	}

	rl, err := readline.New("GCSQL> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("Welcome to GCSQL, the Google Cloud Storage SQL interface.")
	fmt.Println("Type 'EXIT' to quit.")

	for {
		input, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		if strings.ToUpper(input) == "EXIT" {
			fmt.Println("Goodbye!")
			break
		}

		// Add the input to history
		rl.SaveHistory(input)

		// Handle the input
		if err := handleInput(ctx, client, input); err != nil {
			fmt.Printf("Error: %v", err)
		}
	}
}

// func main() {
// 	ctx := context.Background()
// 	client, err := CreateClient(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("Welcome to GCSQL, the Google Cloud Storage SQL interface.")
// 	fmt.Println("Type 'EXIT' to quit.")
// 	for {
// 		fmt.Print("GCSQL> ")
// 		scanner.Scan()
// 		input := scanner.Text()

// 		if strings.ToUpper(input) == "EXIT" {
// 			fmt.Println("Goodbye!")
// 			break
// 		}

// 		// Handle the input
// 		if err := handleInput(ctx, client, input); err != nil {
// 			fmt.Printf("Error: %v", err)
// 		}
// 	}
// }
