import requests
import logging
import sys
import time
import random

status_url = "http://192.168.9.8/include/auth_action.php?k="
login_url = "http://192.168.9.8/include/auth_action.php"
user_name = "2017221102009"
pass_word = "self1098,"


def init_logging():
    # init logger
    logger = logging.getLogger()
    log_file = sys.path[0] + "/log/" + time.strftime("%Y%m%d%H%M%S", time.localtime()) + ".log"
    formatter = "[%(asctime)s] [%(levelname)s] %(message)s"
    file_handler = logging.FileHandler(log_file, encoding="utf-8")

    logging.basicConfig(level="INFO", format=formatter)
    file_handler.setFormatter(logging.Formatter(formatter))
    logger.addHandler(file_handler)


def check_status():
    keys = random.randint(1, 100000)
    payload = {'action': 'get_online_info', 'key': keys}
    res = requests.post(url=status_url+str(keys), data=payload)
    return not (res.text == "not_online")


def login():
    ac_id = 1
    user_mac = ""
    user_ip = ""
    nas_ip = ""
    save_me = 0
    domain = "@uestc"
    ajax = 1
    payload = {
        "action": "login",
        "username": user_name,
        "password": pass_word,
        "ac_id": ac_id,
        "user_mac": user_mac,
        "user_ip": user_ip,
        "nas_ip": nas_ip,
        "save_me": save_me,
        "domain": domain,
        "ajax": ajax,
    }
    res = requests.post(url=login_url, data=payload)
    if "login_ok" in res.text:
        logging.info("login successfully.")
        return True
    else:
        return False


def check_loop():
    while True:
        if not check_status():
            logging.info("not online, now try to login.")
            if not login():
                logging.error("login failed, wait 20 seconds to reconnect.")
                time.sleep(20)
                continue
            else:
                continue
        else:
            logging.info("network online, sleep for 10 minutes.")
            time.sleep(60*10)


def main():
    init_logging()
    logging.info("app started.")
    if not check_status():
        login()
    check_loop()


if __name__ == "__main__":
    main()
