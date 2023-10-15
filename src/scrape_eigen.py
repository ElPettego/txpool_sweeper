#!/usr/bin/env python3.8

import swg_scraper as s
import swg_db as sdb
import static as cfg
import time as t
from selenium.webdriver.common.by import By
from selenium.common.exceptions import NoSuchElementException
import logger as l

DB = sdb.DB('data/db.db')
DB.connect()
SCRAPER = s.Scraper(visible=0)

def main():
    print(l.info(), 'starting scraping eigen')
    txs = []
    SCRAPER.go_to_page(f'{cfg.eigen_base}/mev/bsc/sandwich')
    t.sleep(5)
    for i in range(1, 10):
        tx = SCRAPER.driver.find_element(By.XPATH, f"/html/body/div[1]/div/div/main/div[1]/div/div[4]/section/table/tbody/tr[{i}]/td[2]/div/div/a")
        txs.append(tx.get_attribute('href'))
    # print('txs', txs)
    t.sleep(5)
    
    for tx in txs:
        SCRAPER.go_to_page(tx)
        t.sleep(5)
        for i in range(3, 10):
            try:
                test = SCRAPER.driver.find_element(By.XPATH, f'/html/body/div[1]/div/div/main/div[1]/div/div[3]/section/div/div/div/div/table[1]/thead/tr/th[{i}]/div/div/div/a')
                token_addr = test.get_attribute('href').split('/')[-1]
                # print(token_addr)
                if token_addr not in cfg.SKIP:
                    DB.insert_record_into_table(
                        'tokens',
                        (token_addr,)
                    )

            except NoSuchElementException:
                # print('element not found')
                break

    # t.sleep(100)


if __name__ == "__main__":
    main()