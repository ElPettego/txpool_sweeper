eigen_base: str = "https://eigenphi.io"
coingecko_base: str = "https://api.geckoterminal.com/api/v2"

new_mevs_table: str = "/html/body/div[1]/div/div/main/div[1]/div/div[4]/section/table"
mev_table_addr: str = "/html/body/div[1]/div/div/main/div[1]/div/div[4]/section/table/tbody/tr[1]/td[2]/div/div/a"
tx_addr_bsc:    str = "/html/body/div[1]/div/div/main/div[1]/div/div[1]/div/div/div/div[1]/div/div/div[1]/a[1]"
tk_addr:        str = "/html/body/div[1]/div/div/main/div[1]/div/div[3]/section/div/div/div/div/table[1]/thead/tr" #/th[3]/div/div/div/a"
                    # /html/body/div[1]/div/div/main/div[1]/div/div[3]/section/div/div/div/div/table[1]/thead/tr/th[3]/div/div/div/a

DB_TABLE = 'tokens'

FBNB: str = '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee'
TBNB: str = ''     
WBNB: str = '0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c'
USDT: str = '0x55d398326f99059ff775485246999027b3197955'
GAST: str = '0x0000000000004946c0e9F43F4Dee607b0eF1fA1c'
LINK: str = '0xf8a0bf9cf54bb92f17374d9e9a321e6a111a51bd'

SKIP = [FBNB, TBNB, WBNB, USDT, GAST, LINK]