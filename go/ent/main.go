package main

import (
	"context"
	"fmt"
	"log"

	"github.com/njorda.github.io/webassembly/go/ent/orm_ent"
	"github.com/njorda.github.io/webassembly/go/ent/orm_ent/user"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// client, err := orm_ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client, err := orm_ent.Open("sqlite3", "file:ent?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	ctx := context.Background()
	user1, err := client.User.Create().SetAge(10).SetName("hello").Save(ctx)
	if err != nil {
		panic(err)
	}
	all, err := client.User.Query().All(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(all))

	res, err := client.User.Query().Where(user.ID(user1.ID)).All(ctx)
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		fmt.Println(r)
	}
	fmt.Println("Part two")
	res, err = client.User.Query().Where(user.AgeGT(5)).All(ctx)
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		fmt.Println(r)
	}
}
