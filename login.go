package xenforo_api

type LoginResponse struct {
	Hash string `json:"hash"`
}

func (x *API) Login(username, password string) error {
	callUrl := x.GetCallURL("authenticate")
	q := callUrl.Query()
	q.Set("username", username)
	q.Set("password", password)
	callUrl.RawQuery = q.Encode()

	res := new(LoginResponse)
	if err := x.MakeCall(callUrl, res); err != nil {
		return err
	}

	x.LoginHash = username + ":" + res.Hash
	return nil
}
