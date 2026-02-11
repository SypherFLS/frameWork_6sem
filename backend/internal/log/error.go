package log

type Err struct {
	Id             int    `json:"code"`
	IsDone         bool   `json:"success"`
	Comment        string `json:"message"`
	Identification Iden   `json:"request_id"`
}

type Iden struct {
	Num     int    `json:"num"`
	Process string `json:"process"`
}

func MakeError(id int, com string, num int, proc string) *Err {
	ident := Iden{
		Num:     num,
		Process: proc,
	}
	return &Err{
		Id:             id,
		IsDone:         false,
		Comment:        com,
		Identification: ident,
	}
}

func BackSucsess() *Err {
	return &Err{
		Id:     200,
		IsDone: true,
	}
}

func BackSucsessCreate() *Err {
	return &Err{
		Id:     201,
		IsDone: true,
	}
}
