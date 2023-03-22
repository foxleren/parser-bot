#!/usr/bin/env python
# coding: utf-8

# In[ ]:


import pandas as pd
import numpy as np
import requests
import datetime
import pickle
import os

from bs4 import BeautifulSoup
from time import sleep
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from tqdm import tqdm

# Сброс ограничений на количество выводимых рядов
pd.set_option('display.max_rows', None)

# Сброс ограничений на число столбцов
pd.set_option('display.max_columns', None)


# In[ ]:


def pickle_dump(data, file_name):
    with open(file_name, 'wb') as f:
        pickle.dump(data, f)


def pickle_load(file_name):
    with open(file_name, 'rb') as f:
        data = pickle.load(f)
        return data


def protected_get_html_soup_num(href, class_name='resultsearch_text', num=5):
    response = None
    
#     print(href)
    reload_flag = 0
    while reload_flag < num:
        try:
            headers = {'user-agent': '__________'}
            response = requests.get(href, headers=headers)
            sleep(3)
        except:
            print("Poor connection: reloaded page " + href)
            
            sleep(5)
            reload_flag += 1
            continue
        
        if response is not None and response.status_code == 200:
            raw_html = BeautifulSoup(response.content, features="lxml")
            
            if raw_html.find('div', class_=class_name) is None:
                print("Something wrong with the loaded page " + href)
                
                sleep(5)
                reload_flag += 1
                continue
            
            return raw_html
    
    print('Can not get ' + href)
    return -1


# In[ ]:


court_dict = {
    'mgs': 'Московский городской суд',
    'babushkinskij': 'Бабушкинский районный суд',
    'basmannyj': 'Басманный районный суд',
    'butyrskij': 'Бутырский районный суд',
    'gagarinskij': 'Гагаринский районный суд',
    'golovinskij': 'Головинский районный суд',
    'dorogomilovskij': 'Дорогомиловский районный суд',
    'zamoskvoreckij': 'Замоскворецкий районный суд',
    'zelenogradskij': 'Зеленоградский районный суд',
    'zyuzinskij': 'Зюзинский районный суд',
    'izmajlovskij': 'Измайловский районный суд',
    'koptevskij': 'Коптевский районный суд',
    'kuzminskij': 'Кузьминский районный суд',
    'kuncevskij': 'Кунцевский районный суд',
    'lefortovskij': 'Лефортовский районный суд',
    'lyublinskij': 'Люблинский районный суд',
    'meshchanskij': 'Мещанский районный суд',
    'nagatinskij': 'Нагатинский районный суд',
    'nikulinskij': 'Никулинский районный суд',
    'ostankinskij': 'Останкинский районный суд',
    'perovskij': 'Перовский районный суд',
    'preobrazhenskij': 'Преображенский районный суд',
    'presnenskij': 'Пресненский районный суд',
    'savyolovskij': 'Савёловский районный суд',
    'simonovskij': 'Симоновский районный суд',
    'solncevskij': 'Солнцевский районный суд',
    'taganskij': 'Таганский районный суд',
    'tverskoj': 'Тверской районный суд',
    'timiryazevskij': 'Тимирязевский районный суд',
    'troickij': 'Троицкий районный суд',
    'tushinskij': 'Тушинский районный суд',
    'hamovnicheskij': 'Хамовнический районный суд',
    'horoshevskij': 'Хорошёвский районный суд',
    'cheryomushkinskij': 'Черёмушкинский районный суд',
    'chertanovskij': 'Чертановский районный суд',
    'shcherbinskij': 'Щербинский районный суд',
}

possible_names = ['скилбокс', 'скиллбокс', 'skillbox']


# In[ ]:


res_num = pickle_load('res_num.pkl')
res_num


# In[ ]:


res_data = pickle_load('res_data.pkl')
res_data.head()


# In[ ]:


def get_court_data(name='скилбокс'):
    res_data = dict()
    res_data['Запрос'] = []
    res_data['Суд'] = []
    res_data['Ссылка'] = []
    res_data['Номер дела'] = []
    # res_data['Дата парсинга'] = []
    
    href_prefix = "https://mos-gorsud.ru/search?formType=shortForm&participant=" + name
    
    for court in tqdm(court_dict.keys()):
        href = href_prefix + '&courtAlias=' + court
        raw_html = protected_get_html_soup_num(href)
        
        if raw_html == -1:
            return -1
        
        resultsearch_text = raw_html.find('div', class_='resultsearch_text')
        resultsearch_text = resultsearch_text.text.strip().lower()
        
        if 'ничего не найдено' not in resultsearch_text:
            resultsearch_text = resultsearch_text.split('\n')[0].strip()
            amount_found = int(resultsearch_text.split(':')[1].strip())
            
            table = raw_html.find('table', class_='custom_table').find('tbody')
            table_rows = table.find_all('tr')
            
            for row in table_rows:
                request_res = name
                court_res = court_dict[court]
                href_res = 'https://mos-gorsud.ru' + row.find('nobr').find('a').get('href')
                nobr_res = row.find('nobr').text.strip()

                res_data['Запрос'].append(request_res)
                res_data['Суд'].append(court_res)
                res_data['Ссылка'].append(href_res)
                res_data['Номер дела'].append(nobr_res)
            
            if len(table_rows) != amount_found:
                new_href = href + '&page='
                page_num = 2
                
                raw_html_new = protected_get_html_soup_num(new_href + str(page_num))
                
                if raw_html_new == -1:
                    return -1
                
                resultsearch_text = raw_html_new.find('div', class_='resultsearch_text')
                resultsearch_text = resultsearch_text.text.strip().lower()
                
                while 'ничего не найдено' not in resultsearch_text:
                    table = raw_html_new.find('table', class_='custom_table').find('tbody')
                    table_rows = table.find_all('tr')

                    for row in table_rows:
                        request_res = name
                        court_res = court_dict[court]
                        href_res = 'https://mos-gorsud.ru' + row.find('nobr').find('a').get('href')
                        nobr_res = row.find('nobr').text.strip()

                        res_data['Запрос'].append(request_res)
                        res_data['Суд'].append(court_res)
                        res_data['Ссылка'].append(href_res)
                        res_data['Номер дела'].append(nobr_res)
                    
                    page_num += 1
                    
                    raw_html_new = protected_get_html_soup_num(new_href + str(page_num))
                    
                    if raw_html_new == -1:
                        return -1
                    
                    resultsearch_text = raw_html_new.find('div', class_='resultsearch_text')
                    resultsearch_text = resultsearch_text.text.strip().lower()
    
    return pd.DataFrame(res_data)


# In[ ]:


res_data1 = get_court_data()
res_data1.head()


# In[ ]:


# res_data2 = get_court_data('skilbox')
# res_data2.head()


# In[ ]:


# res_data = pd.concat([res_data, res_data2], ignore_index=True)
# res_data.head()


# In[ ]:


pickle_dump(res_data, 'res_data.pkl')


# In[ ]:


# res_data.to_excel('court_data.xlsx')


# In[ ]:


# len(res_data)


# In[ ]:


file_path = 'court_data.xlsx'


# In[ ]:


global prev_df
new_lines_counter = 0


# In[ ]:


def highlight_new_rows(row):
    if not prev_df.empty and any([val not in prev_df.values for val in row.values]):
        return ['background-color: #9EC3EF'] * len(row)
    else:
        return [''] * len(row)
def highlight_all_rows(row):
    return ['background-color: #9EC3EF'] * len(row)


# In[ ]:


if os.path.isfile(file_path):
    prev_df = pd.read_excel('court_data.xlsx')
    diff = pd.concat([prev_df, res_data]).drop_duplicates(keep=False)
    new_lines_counter = len(diff)
    merged_table = pd.merge(prev_df,res_data, how='outer')
    merged_table.fillna(value='', inplace=True)
    merged_table.style.apply(highlight_new_rows, axis=1).to_excel(file_path, engine='openpyxl', index=False)
else:
    new_lines_counter = len(res_data)
    res_data.style.apply(highlight_all_rows, axis=1).to_excel('court_data.xlsx', index=False)


# In[ ]:


import subprocess


# In[ ]:


text = str(new_lines_counter)
subprocess.run(["echo", text])

