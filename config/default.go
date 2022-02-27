package config

const defaultJson string = `{
	"service": {
		"address": ":7074",
		"ttl": 15,
		"interval": 10
	},
	"logger": {
		"level": "info",
		"file": "logs/server.log",
		"std": false
	},
	"database": {
		"name": "rgsCloud",
		"ip": "192.168.1.31",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	}
}
`
