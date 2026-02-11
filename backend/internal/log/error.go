package log 


import (

)

type Error struct {
	id int 
	isDone bool
	comment string 
	identification Iden
}

type Iden struct {
	num int 
	process string
}


func MakeError(id int, com string, num int, proc string) *Error{
	ident := Iden {
		num:num, 
		process: proc,
	}
	return &Error {
		id:id,
		isDone: false,
		comment:com,
		identification: ident,
	}
}

func BackSucsess() *Error{
	return &Error{
		isDone: true,
	}
}