CREATE TABLE card (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    bank TEXT NOT NULL,
    credit_limit NUMERIC(10,2) NOT NULL CHECK (credit_limit >= 0),
    due_day SMALLINT NOT NULL CHECK (due_day BETWEEN 1 AND 28),

    CONSTRAINT bank_check
        CHECK (bank IN ('nubank'))
);

CREATE TABLE fatura (
    id TEXT PRIMARY KEY,
    card_id TEXT NOT NULL,
    year SMALLINT NOT NULL CHECK (year >= 2000),
    month SMALLINT NOT NULL CHECK (month BETWEEN 1 AND 12),
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    total NUMERIC(10,2) NOT NULL,

    CONSTRAINT fk_card
        FOREIGN KEY (card_id)
        REFERENCES card(id)
        ON DELETE CASCADE,

    CONSTRAINT status_check
        CHECK (status IN ('pending', 'paid')),

    CONSTRAINT unique_card_year_month
        UNIQUE (card_id, year, month)
);

CREATE TABLE bill (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    date TEXT NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    fatura_id TEXT NOT NULL,
    method TEXT NOT NULL,
    category TEXT NOT NULL,

    CONSTRAINT fk_fatura
        FOREIGN KEY (fatura_id)
        REFERENCES fatura(id)
        ON DELETE CASCADE,

    CONSTRAINT method_check
        CHECK (method IN ('parcelado', 'fixo')),

    CONSTRAINT category_check
        CHECK (category IN (
            'transporte',
            'alimentação',
            'mercado',
            'saúde',
            'assinaturas',
            'vestuario',
            'celular',
            'entretenimento',
            'varejo',
            'moradia',
            'educacao',
            'viagem',
            'servicos',
            'outros'
        ))
);