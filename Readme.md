# Desafio Abertura e fechamento do Leilão - Go Routines

Objetivo: Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

# Subir a aplicação localmente/Rodar o Dockerfile

Executar na pasta root no cmd o seguinte comando: `docker compose up -d`.

# Validação da aplicação

Ao criar o leilão, é validado de segundo em segundo o tempo de expiração do leilão, caso o tempo atual ultrapasse o tempo estipulado pelo leilão, então deve marcar o leilão com enum status Completed (1). Para validar só olhar o banco de dados do mongo DB gerado. (Olhar arquivo .env)


