package xenforo_api

import (
	"sort"
	"strconv"
	"strings"
)

type UserResponse struct {
	UserID   int    `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Title    string `json:"custom_title,omitempty"`

	PrimaryGroup_    int    `json:"user_gorup_id,omitempty"`       // An internal field. This is pushed onto the beginning of the groups array and then set to zero.
	SecondaryGroups_ string `json:"secondary_group_ids,omitempty"` // An internal field. The secondary group IDs are transmitted as a string rather than an int, so we have to decode them our selves.
	Groups           []int  `json:"-"`
}

func (u *UserResponse) IsInGroup(group int) bool {
	for _, groupID := range u.Groups {
		if group == groupID {
			return true
		}
	}

	return false
}

func (u *UserResponse) Initialize() {
	secondaryGroupsSplit := strings.Split(u.SecondaryGroups_, ",")
	u.Groups = make([]int, 1, len(secondaryGroupsSplit)+1)
	u.Groups[0] = u.PrimaryGroup_
	for _, groupStr := range secondaryGroupsSplit {
		groupStr = strings.TrimSpace(groupStr)
		if len(groupStr) > 0 {
			groupID, err := strconv.ParseInt(groupStr, 0, 32)
			if err != nil {
				panic(err)
			}
			u.Groups = append(u.Groups, int(groupID))
		}
	}

	sort.Ints(u.Groups)
}

func (x *API) GetUser(id string) (*UserResponse, error) {
	callUrl := x.GetCallURL("getUser")
	q := callUrl.Query()
	if len(id) > 0 {
		q.Set("value", id)
	} else {
		q.Del("value")
	}
	callUrl.RawQuery = q.Encode()

	res := new(UserResponse)
	if err := x.MakeCall(callUrl, res); err != nil {
		return nil, err
	}
	res.Initialize()

	return res, nil
}
