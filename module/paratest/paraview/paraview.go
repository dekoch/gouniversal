package paraview

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/dekoch/gouniversal/module/paratest/auftrag"
	"github.com/dekoch/gouniversal/module/paratest/einzelpara"
	"github.com/dekoch/gouniversal/module/paratest/parameter"
	"github.com/google/uuid"
)

const (
	PARACNT        = 5
	TYPPARAMETER   = "Typ"
	IDENTPARAMETER = "Ident"
)

type ParaView struct {
	Auftrag auftrag.Auftrag
	Typ     []ParaSet
	Ident   []ParaSet
}

type ParaSet struct {
	Header    parameter.Parameter
	Parameter einzelpara.EinzelPara
}

func (pv *ParaView) Save(db *sql.DB, tx *sql.Tx) error {

	var err error

	func() {

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				tx, err = db.Begin()

			case 1:
				// save Auftrag
				err = pv.Auftrag.Save(tx)

			case 2:
				// save Typ Parameter
				for _, p := range pv.Typ {

					err = p.Header.Save(tx)
					if err != nil {
						return
					}

					err = p.Parameter.Save(tx)
					if err != nil {
						return
					}
				}

			case 3:
				// save Ident Parameter
				for _, p := range pv.Ident {

					err = p.Header.Save(tx)
					if err != nil {
						return
					}

					err = p.Parameter.Save(tx)
					if err != nil {
						return
					}
				}

			case 4:
				err = tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	if err != nil {
		tx.Rollback()
	}

	return err
}

func (pv *ParaView) Load(auftr, prozess string, db *sql.DB) error {

	var (
		err   error
		found bool
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				// load Auftrag
				found, err = pv.Auftrag.Load(auftr, db)
				if found == false {
					err = errors.New("not found")
					return
				}

			case 1:
				// load Typ Parameter UUIDs
				for n := 1; n <= PARACNT; n++ {

					var ps ParaSet
					found, err = ps.Header.Load(pv.Auftrag.TypUUID, prozess, strconv.Itoa(n), db)
					if err != nil {
						return
					}

					if found {
						pv.Typ = append(pv.Typ, ps)
					}
				}

			case 2:
				// load Typ Parameters
				for i, p := range pv.Typ {

					var ep einzelpara.EinzelPara
					found, err = ep.Load(p.Header.ParameterUUID, db)
					if err != nil {
						return
					}

					if found {
						pv.Typ[i].Parameter = ep
					}
				}

			case 3:
				// load Ident Parameter UUIDs
				for n := 1; n <= PARACNT; n++ {

					var ps ParaSet
					found, err = ps.Header.Load(pv.Auftrag.IdentUUID, prozess, strconv.Itoa(n), db)
					if err != nil {
						return
					}

					if found {
						pv.Ident = append(pv.Ident, ps)
					}
				}

			case 4:
				// load Ident Parameters
				for i, p := range pv.Ident {

					var ep einzelpara.EinzelPara
					found, err = ep.Load(p.Header.ParameterUUID, db)
					if err != nil {
						return
					}

					if found {
						pv.Ident[i].Parameter = ep
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (pv *ParaView) NewAuftrag(name string) string {

	return pv.Auftrag.New(name)
}

func (pv *ParaView) NewParameterSet(isTypParameter bool, prozess, name string) {

	var list []ParaSet

	u := uuid.Must(uuid.NewRandom())
	paraUUID := u.String()

	for i := 1; i <= PARACNT; i++ {

		var ps ParaSet
		ps.Header.UUID = paraUUID

		if isTypParameter {
			ps.Header.Typ = TYPPARAMETER
		} else {
			ps.Header.Typ = IDENTPARAMETER
		}

		ps.Header.Name = name
		ps.Header.Prozess = prozess
		ps.Header.ParameterNr = strconv.Itoa(i)

		u = uuid.Must(uuid.NewRandom())
		ps.Header.ParameterUUID = u.String()

		list = append(list, ps)
	}

	for i, ps := range list {

		var ep einzelpara.EinzelPara
		ep.UUID = ps.Header.ParameterUUID

		if isTypParameter {
			ep.Name = name + " Typ Parameter " + ps.Header.ParameterNr
		} else {
			ep.Name = name + " Ident Parameter " + ps.Header.ParameterNr
		}

		ep.Wert = ps.Header.ParameterNr

		list[i].Parameter = ep
	}

	if isTypParameter {
		pv.Auftrag.TypUUID = paraUUID
		pv.Typ = list
	} else {
		pv.Auftrag.IdentUUID = paraUUID
		pv.Ident = list
	}
}
