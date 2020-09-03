package controllers

import (
	STRUCTURES "../structures"
)

//Variable global que almacenara al usiario logueado
var UserLogued = STRUCTURES.USER{User_isLoged: false}

//Inserta el usuario en la estructura
func AddLogedUser(id string, types string, group string, username string, password string) {
	UserLogued = STRUCTURES.USER{
		User_id:       id,
		User_type:     types,
		User_group:    group,
		User_username: username,
		User_password: password,
		User_isLoged:  true}
}

//Limpia la sesion
func AddLogOut() (bool, string) {
	if IsLogged() {
		var temp = UserLogued.User_username
		UserLogued = STRUCTURES.USER{User_isLoged: false}
		return true, temp
	}
	return false, ""
}

//Verifica si esta o no loguqo un usuario
func IsLogged() bool {
	if UserLogued.User_isLoged {
		return true
	}
	return false
}

//Verifica si esta o no loguqo un ROOT logueado
func IsRootLogged() bool {
	if UserLogued.User_isLoged && (UserLogued.User_group == "root") {
		return true
	}
	return false
}

//Retorna el usuario logueado
func GetLogedUser() STRUCTURES.USER {
	if IsLogged() {
		return UserLogued
	}
	return STRUCTURES.USER{}
}
