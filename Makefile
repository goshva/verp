.PHONY: run dev build clean

# Запуск с горячей перезагрузкой (air)
dev:
	air

# Обычный запуск без перезагрузки
run:
	go run cmd/server/main.go

# Сборка бинарника
build:
	go build -o bin/server cmd/server/main.go

# Очистка
clean:
	rm -rf tmp/ bin/

# Установка air
install-air:
	go install github.com/air-verse/air@latest

# Проверка установки air
check-air:
	air -v