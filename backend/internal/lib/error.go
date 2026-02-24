package lib

type Err struct {
	Id int `json:"code"`
	IsDone bool `json:"success"`
	Comment string `json:"message"`
	Identification Iden `json:"request_id"`
}

// type CallHttp struct {
// 	Method string
// 	Status int
// 	Body any
// }

type Iden struct {
	Num int `json:"num"`
	Process string `json:"process"`
}



func MakeError(id int, com string, num int, proc string) *Err {
	ident := Iden{
		Num:     num,
		Process: proc,
	}
	err := &Err{
		Id: id,
		IsDone: false,
		Comment: com,
		Identification: ident,
	}
	LogError(err)
	return err
}

func BackSucsess(id int, com string, num int, proc string) *Err {
	ident := Iden{
		Num: num,
		Process: proc,
	}
	err := &Err{
		Id: id,
		IsDone: true,
		Comment: com,
		Identification: ident,
	}
	LogSuccess(err)
	return err
}
