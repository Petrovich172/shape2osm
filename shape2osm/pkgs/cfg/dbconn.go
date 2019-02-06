package cfg

import (
	"github.com/go-pg/pg"
)

// Connections - подключения к базам данных Postgres
type Connections map[string]*pg.DB

// EstablishAllConnections - подключение ко всем указанным в JSON базам данных Postgres
func (cons *Connections) EstablishAllConnections(cfg *PostgresDatabaseCfg) {
	(*cons) = make(map[string]*pg.DB)
	for i := range *cfg {
		if i == 1 { // Skip telegramm connection
			continue
		} //
		(*cons)[(*cfg)[i].Database] = pg.Connect(&pg.Options{
			Addr:      (*cfg)[i].Host + ":" + (*cfg)[i].Port,
			User:      (*cfg)[i].User,
			Password:  (*cfg)[i].Password,
			Database:  (*cfg)[i].Database,
			TLSConfig: (*cfg)[i].EnableTLS,
		})
	}
}
