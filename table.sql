CREATE TABLE "Draws" (
	"drawDate"	TEXT,
	"ID"	INTEGER NOT NULL UNIQUE,
	"drawTime"	TEXT NOT NULL,
	"extraPick"	INTEGER NOT NULL,
	PRIMARY KEY("ID")
);
CREATE TABLE "Picks" (
	"drawID"	INTEGER NOT NULL,
	"pickNum"	INTEGER NOT NULL
);
CREATE UNIQUE INDEX "draw_pick_asc" ON "Picks" (
	"drawID"	ASC,
	"pickNum"	ASC
);
CREATE UNIQUE INDEX "draw_pick_desc" ON "Picks" (
	"drawID"	DESC,
	"pickNum"	ASC
);
CREATE UNIQUE INDEX "draw_date_asc" ON "Draws" (
	"ID"	ASC,
	"drawDate"	ASC
);