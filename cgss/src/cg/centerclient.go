package cg

import (
    "json"
)

type CenterClient struct {
    *ipc.IpcClient
}

func (client *CenterClient)addPlayer(player *Player) error {
    b, err := json.Marshal(*player);
    if err != nil {
        return err
    }

    resp, err := client.Call("addplayer", string(b))
    if err == nil && resp.Code == "200" {
        return nil
    }
    return err
}

func (client *CenterClient)RemovePlayer(name string) error {
    ret, _ := client.Call("listplayer", name)
    if ret.Code == "200" {
        return nil
    }
    return errors.New(ret.Code)
}

func (client *CenterClient)ListPlayer(params string) (ps []*player, err error) {
    resp, _ := client.Call("listplayer", params)
    if resp.Code != "200" {
        err = errors.New(resp.Code)
        return
    }

    err = json.Unmarshal([]byte(resp.Body), &ps)
    return
}
