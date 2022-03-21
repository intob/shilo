build:
	@printf "\033[32m\xE2\x9c\x93Building shilo with tag shilo:latest\n\033[0m"
	@docker build -t shilo:latest .

run:
	@printf "\033[32m\xE2\x9c\x93Running shilo:latest\n\033[0m"
	@docker run -d -p 9000:9000 shilo:latest