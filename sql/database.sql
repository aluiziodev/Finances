CREATE TABLE bill (
    id TEXT PRIMARY KEY, 
    title TEXT NOT NULL,
    date TEXT NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    fatura TEXT NOT NULL,
    method TEXT NOT NULL,
    category TEXT NOT NULL,

    CONSTRAINT fk_fatura
        FOREIGN KEY (fatura)
        REFERENCES fatura(id)
        ON DELETE CASCADE,
    
    CONSTRAINT method_check
        CHECK (method in ('parcelado', 'fixo')),
    
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



CREATE TABLE fatura (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    bank TEXT NOT NULL,
    status TEXT NOT NULL,
    total NUMERIC(10,2) NOT NULL,

    CONSTRAINT status_check
        CHECK (status IN ('pending', 'paid')),

    CONSTRAINT bank_check
        CHECK (bank IN ('nubank'))
)