package repositories

import "github.com/jyotirmoydotdev/openfy/db/models"

var Admins []models.Admin
var AdminSecrets = make(map[string]string)
var AdminIDCounter int
