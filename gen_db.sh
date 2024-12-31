gentool \
    -db mysql \
    -dsn "XXX:XXX@tcp(localhost:3306)/XXX?charset=utf8mb4&parseTime=True&loc=Local" \
    -onlyModel\
    -fieldNullable \
    -modelPkgName "model" \
    -outPath "./internal/domain/model"