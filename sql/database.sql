CREATE TABLE bill (
    id TEXT PRIMARY KEY, 
    title TEXT NOT NULL,
    date TEXT NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    fatura TEXT NOT NULL,

    CONSTRAINT fk_fatura
        FOREIGN KEY (fatura)
        REFERENCES fatura(id)
        ON DELETE CASCADE
);



CREATE TABLE fatura (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    total NUMERIC(10,2) NOT NULL,

    CONSTRAINT status_check
        CHECK (status IN ('pending', 'paid'))
)