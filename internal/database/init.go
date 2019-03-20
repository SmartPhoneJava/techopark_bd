package database

import (
	"database/sql"
	"escapade/internal/config"
	"fmt"
	"os"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {

	// for local launch
	if os.Getenv(CDB.URL) == "" {
		os.Setenv(CDB.URL, "user=rolepade password=escapade dbname=escabase sslmode=disable")
	}

	var database *sql.DB
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		fmt.Println("database/Init cant open:" + err.Error())
		return
	}

	db = &DataBase{
		Db: database,
	}
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = db.Db.Ping(); err != nil {
		fmt.Println("database/Init cant access:" + err.Error())
		return
	}
	fmt.Println("database/Init open")
	//if !db.areTablesCreated(CDB.Tables) {
	if err = db.CreateTables(); err != nil {
		return
	}
	//}

	return
}

func (db *DataBase) checkTable(tableName string) (err error) {
	sqlStatement := `
    SELECT count(1)
  FROM information_schema.tables tbl 
  where tbl.table_name like $1;`
	row := db.Db.QueryRow(sqlStatement, tableName)

	var result int
	if err = row.Scan(&result); err != nil {
		fmt.Println(tableName + " doesnt exists. Create it!" + err.Error())

		return
	}
	return
}

func (db *DataBase) areTablesCreated(tables []string) (created bool) {
	created = true
	for _, table := range tables {
		if err := db.checkTable(table); err != nil {
			created = false
			break
		}
	}
	return
}

func (db *DataBase) CreateTables() error {
	sqlStatement := `
	DROP TABLE IF EXISTS Session;
    DROP TABLE IF EXISTS Game;
    DROP TABLE IF EXISTS Player;
    DROP TABLE IF EXISTS Photo;

    DROP TABLE IF EXISTS Post;
    DROP TABLE IF EXISTS Thread;
    DROP TABLE IF EXISTS Forum;
    DROP TABLE IF EXISTS UserForum;

    CREATE Table UserForum (
        id SERIAL PRIMARY KEY,
        nickname varchar(80) UNIQUE NOT NULL,
        fullname varchar(30) NOT NULL,
        email varchar(50) UNIQUE NOT NULL,
        about varchar(1000) 
    );

    CREATE Table Forum (
        id SERIAL PRIMARY KEY,
        posts int default 0,
        slug varchar(80) not null UNIQUE,
        threads int default 0,
        title varchar(120) not null,
        user_nickname varchar(80) not null
    );

    ALTER TABLE Forum
        ADD CONSTRAINT forum_user
        FOREIGN KEY (user_nickname)
        REFERENCES UserForum(nickname)
            ON DELETE CASCADE;

    CREATE Table Thread (
        id SERIAL PRIMARY KEY,
        author varchar(120) not null,
        forum varchar(120) not null,
        message varchar(1600) not null,
        created    TIMESTAMPTZ,
        title varchar(120) not null,
        votes int default 0,
        slug varchar(120) default null
    );

    ALTER TABLE Thread
    ADD CONSTRAINT thread_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    ALTER TABLE Thread
    ADD CONSTRAINT thread_forum
    FOREIGN KEY (forum)
    REFERENCES Forum(slug)
        ON DELETE CASCADE;

    CREATE Table Post (
        id SERIAL PRIMARY KEY,
        author varchar(120) not null,
        forum varchar(120),
        message varchar(1600) not null,
        created    TIMESTAMPTZ,
        isEdited boolean default false,
        thread int ,
        parent int
    );

    ALTER TABLE Post
    ADD CONSTRAINT post_user
    FOREIGN KEY (author)
    REFERENCES UserForum(nickname)
        ON DELETE CASCADE;

    ALTER TABLE Post
    ADD CONSTRAINT post_forum
    FOREIGN KEY (forum)
    REFERENCES Forum(slug)
        ON DELETE CASCADE;
    
    ALTER TABLE Post
    ADD CONSTRAINT post_thread
    FOREIGN KEY (thread)
    REFERENCES Thread(id)
        ON DELETE CASCADE;

	CREATE TABLE Player (
    id SERIAL PRIMARY KEY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_title varchar(50),
    --FirstSeen   timestamp without time zone NOT NULL,
	--LastSeen    timestamp without time zone NOT NULL,
    best_score  int default 0 CHECK (best_score > -1),
    best_time   int default 0 CHECK (best_time > -1),
    GamesTotal  int default 0 CHECK (GamesTotal > -1),
	SingleTotal int default 0 CHECK (SingleTotal > -1),
	OnlineTotal int default 0 CHECK (OnlineTotal > -1),
	SingleWin   int default 0 CHECK (SingleWin > -1),
	OnlineWin   int default 0 CHECK (OnlineWin > -1),
	MinsFound   int default 0 CHECK (MinsFound > -1)
    
);

CREATE Table Session (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    session_code varchar(30) NOT NULL,
    expiration timestamp without time zone NOT NULL
);

ALTER TABLE Session
ADD CONSTRAINT session_player
   FOREIGN KEY (player_id)
   REFERENCES Player(id)
   ON DELETE CASCADE;

CREATE Table Game (
    id SERIAL PRIMARY KEY,
    player_id   int NOT NULL,
    FieldWidth  int CHECK (FieldWidth > -1),
    FieldHeight int CHECK (FieldHeight > -1),
    MinsTotal   int CHECK (MinsTotal > -1),
    MinsFound   int CHECK (MinsFound > -1),
    Finished bool NOT NULL,
    Exploded bool NOT NULL,
    Date timestamp without time zone NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Player (id)
);

--GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO escapade;

INSERT INTO Player(name, password, email, best_score, best_time) VALUES
    ('tiger', 'Bananas', 'tinan@mail.ru', 1000, 10),
    ('panda', 'apple', 'today@mail.ru', 2323, 20),
    ('catmate', 'juice', 'allday@mail.ru', 10000, 5),
    ('hotdog', 'where', 'three@mail.ru', 88, 1000);

    /*
    id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_id int,
    best_score int,
    FOREIGN KEY (photo_id) REFERENCES Photo (id)
    */

INSERT INTO Game(player_id, FieldWidth, FieldHeight,
MinsTotal, MinsFound, Finished, Exploded, Date) VALUES
    (1, 50, 50, 100, 20, true, true, date '2001-09-28'),
    (1, 50, 50, 80, 30, false, false, date '2018-09-27'),
    (1, 50, 50, 70, 70, true, false, date '2018-09-26'),
    (1, 50, 50, 60, 30, true, true, date '2018-09-23'),
    (1, 50, 50, 50, 50, true, false, date '2018-09-24'),
    (1, 50, 50, 40, 30, true, false, date '2018-09-25'),
    (2, 25, 25, 80, 30, false, false, date '2018-08-27'),
    (2, 25, 25, 70, 70, true, false, date '2018-08-26'),
    (2, 25, 25, 60, 30, true, true, date '2018-08-23'),
    (3, 10, 10, 10, 10, true, false, date '2018-10-26'),
    (3, 10, 10, 20, 19, true, true, date '2018-10-23'),
    (3, 10, 10, 30, 30, true, false, date '2018-10-24'),
    (3, 10, 10, 40, 5, true, false, date '2018-10-25');

    /*
CREATE Table Game (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    FieldWidth int CHECK (FieldWidth > -1),
    FieldHeight int CHECK (FieldHeight > -1),
    MinsTotal int CHECK (MinsTotal > -1),
    MinsFound int CHECK (MinsFound > -1),
    Finished bool NOT NULL,
    Exploded bool NOT NULL,
    Date timestamp without time zone NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Player (id)
);
    */
	`
	_, err := db.Db.Exec(sqlStatement)

	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}
	return err
}
