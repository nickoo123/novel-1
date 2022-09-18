from selenium import webdriver
import time
option = webdriver.ChromeOptions()
option.add_argument('headless')
option.add_argument('no-sandbox')
option.add_argument('disable-dev-shm-usage')

def getHtml(url, loadmore = False, waittime = 2):
    browser = webdriver.Chrome('/app/chromedriver',chrome_options=option)
    browser.get(url)
    if loadmore:
        while True:
            try:
                time.sleep(waittime)
                browser.find_element_by_name("Email").send_keys("william88456@gmail.com")
#                 browser.find_element_by_id("Passwd-hidden").send_keys("3rdVRhD5DrZb")
                next_cont = browser.find_element_by_id("next").click()
                print(browser.page_source)
            except:
                break
    html = browser.page_source
    browser.quit()
    return html

getHtml('https://accounts.google.com/o/oauth2/v2/auth?scope=https://mail.google.com&access_type=offline&include_granted_scopes=true&state=state_parameter_passthrough_value&redirect_uri=https://www.biqugesk.cc&response_type=code&client_id=423816617029-k0jrtlbr362qn16ca82tcjbatqv53v4r.apps.googleusercontent.com',True,30)
# print(browser.page_source)
# elements = browser.find_element_by_class_name("browser")
# print(elements)


