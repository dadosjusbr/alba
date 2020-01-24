# Alba
## Sistema de Orquestração de Execuções DadosJusBR

Sistema para orquestração de atividades visando a automatização de processos do DadosJusBR. A função principal da ferramenta é gerenciar o processo de libertação de dados, que inclui as etapas de:

- **Coleta:** Etapa responsável por encontrar, fazer o download dos arquivos e consolidar/traduzir as informações para um formato único do DadosJusBr. Cada coletor é responsável por um determinado órgão e recebe como parâmetro um mês/ano e o identificador do órgão
- **Validação:** Responsável por fazer validações nos dados de acordo a cada contexto;
- **Empacotamento:** Responsável por padronizar os dados no formato de datapackages;
- **Armazenamento:** Responsável por armazenar os dados extraídos, além de versionar também os artefatos baixados e gerados durante a coleta; 

