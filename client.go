package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

const READ = true
const WRITE = false

// The number of Main blocks to Number of Holding Blocks ratio

const nPlusM = 9
var MainToHoldingRatio = 2
var N = int( (nPlusM * MainToHoldingRatio) / (1 + MainToHoldingRatio))
var M = int(nPlusM - N)
var posMap = make(map[int]int)

// Data is the structure of our data. Server does not have this view. For now kept it simple.
type Data struct {
	D int
}

type Req struct{
	Mode bool
	D    Data
	Meta int
}

// Simply calls the cloud read function in this setup.

func myCloudReader(index int)(Data,error){
	if index >= nPlusM {
		return Data{0},errors.New("accessing beyond nPlusM")
	}
	var retData Data
	x,err:= cloudRead(index)
	if err!=nil{return Data{0},errors.New("error while reading from server")}
	x = decrypt(x, "mrigankaismyname")
	err = json.Unmarshal(x,&retData)
	if err!= nil{ return Data{0}, errors.New("unmarshalling error")}
	return retData, nil
}

// The implementation of DetWoORAM
// Algorithm 2 of : https://arxiv.org/pdf/1706.03827.pdf
//
// Assuming that the adversary is a remote user, who can only take snapshot of our cloud storage,
// and can not see our read requests. Since he can take snapshots, he gets to see what we write.
// Thus we need write obliviousness, and not read obliviousness.

func myCloudWriter(meta int, D Data, count int) error{
	if meta >= N { return errors.New("accessing beyond set parameters")}
	encD,err := json.Marshal(D)
	if err!=nil {return err}
	encD = encrypt(encD,"mrigankaismyname")
	err = cloudWrite(N + (count % M), encD)
	if err!=nil{return err	}
	posMap[meta]= N + (count % M)
	start:= ( count * (N / M) ) % N
	end:= ( (count + 1) * (N / M) ) % N
	if end==0{ end = N}
	for i := start; i < end; i++{
		x, B := posMap[i]
		if B== false{
			x = i
		}
		Bytes, err:= cloudRead(x)
		if err!=nil{return err}
		freshBytes := encrypt(decrypt(Bytes,"mrigankaismyname"),"mrigankaismyname")
		posMap[i]=i
		err = cloudWrite(i,freshBytes)
		if err!=nil{return err}
	}
	return nil
}

// Just initiates the server, as in it sets garbage to all locations in the server.
func initiateCloud(){
	x := encrypt([]byte("0"),"mrigankaismyname")
	for i:=0;i< nPlusM;i++ {
		cloudWrite(i,x)
		posMap[i]=i
	}
}

func main(){
	initiateCloud()
	var accessPattern = [14]Req{
		{WRITE, Data{0},0},{WRITE, Data{1},1},{WRITE, Data{2},2},
		{WRITE, Data{3},3},{WRITE, Data{4},4}, {WRITE, Data{5},5},
		{WRITE, Data{6},0}, {WRITE, Data{7},1},

		{READ, Data{0},0},{READ, Data{0},1},	{READ, Data{0},2},
		{READ, Data{0},3}, {READ, Data{0},4}, {READ, Data{0},5},
	}
	var count=0
	for i := range accessPattern{
		if accessPattern[i].Mode == READ {
			x,err :=myCloudReader(posMap[accessPattern[i].Meta])
			if err==nil{
				fmt.Println(x)
			} else {
				fmt.Println(err)
			}

		} else {
			err:= myCloudWriter(accessPattern[i].Meta, accessPattern[i].D, count)
			count= count+1
			if err== nil{
				fmt.Println("Successfully written")
			} else {
				fmt.Println(err)
			}
		}
	}
}
