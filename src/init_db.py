#!/usr/bin/env python3.8

import swg_db as db

DB = db.DB('data/db.db')


def main():
    DB.connect()
    DB.create_table(
        "tokens",
        {
            'address': 'string_pk',
        }
    )
    DB.create_table(
        "sent",
        {
            'address': 'string_pk',
            'creation': 'string',
            'name': 'string'
        }
    )
    
    # DB.insert_record_into_table("tokens", ('0xlkoljos', '2023/25/10'))

if __name__ == "__main__":
    main()