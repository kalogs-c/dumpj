---------------- Natureza Juridica ----------------

-- name: CreateNaturezaJuridica :exec
INSERT INTO naturezas_juridicas (descricao, codigo)
VALUES (?, ?)
ON CONFLICT (codigo) DO UPDATE SET descricao = EXCLUDED.descricao;

---------------- Municipio ----------------

-- name: CreateMunicipio :exec
INSERT INTO municipios (codigo, descricao)
VALUES (?, ?)
ON CONFLICT (codigo) DO UPDATE SET descricao = EXCLUDED.descricao;

---------------- Cnae ----------------

-- name: CreateCnae :exec
INSERT INTO cnaes (codigo, descricao)
VALUES (?, ?)
ON CONFLICT (codigo) DO UPDATE SET descricao = EXCLUDED.descricao;

---------------- Empresa ----------------

-- name: CreateEmpresa :exec
INSERT INTO empresas (cnpj_basico, razao_social, natureza_juridica, capital_social, porte_empresa)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (cnpj_basico) DO UPDATE 
SET 
    razao_social = EXCLUDED.razao_social,
    natureza_juridica = EXCLUDED.natureza_juridica,
    capital_social = EXCLUDED.capital_social,
    porte_empresa = EXCLUDED.porte_empresa;

---------------- Estabelecimento ----------------

-- name: CreateEstabelecimento :exec
INSERT INTO estabelecimentos (
    cnpj_basico, cnpj_ordem, cnpj_dv, identificador_matriz_filial, 
    nome_fantasia, data_abertura, cnae, logradouro, 
    numero, complemento, bairro, cep,
    uf, municipio, ddd, telefone, email
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (cnpj_basico, cnpj_ordem, cnpj_dv) DO UPDATE 
SET 
    nome_fantasia = EXCLUDED.nome_fantasia,
    data_abertura = EXCLUDED.data_abertura,
    cnae = EXCLUDED.cnae,
    logradouro = EXCLUDED.logradouro,
    numero = EXCLUDED.numero,
    complemento = EXCLUDED.complemento,
    bairro = EXCLUDED.bairro,
    cep = EXCLUDED.cep,
    uf = EXCLUDED.uf,
    municipio = EXCLUDED.municipio,
    ddd = EXCLUDED.ddd,
    telefone = EXCLUDED.telefone,
    email = EXCLUDED.email;


