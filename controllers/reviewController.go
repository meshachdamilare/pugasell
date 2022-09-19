package controllers

import (
	"github.com/Christomesh/pugasell/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var ReviewCollection *mongo.Collection = db.OpenCollection(db.Client, "review")
