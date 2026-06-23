import requests
import json
import re
from bs4 import BeautifulSoup
import time

def get_gun_data(url):
    url = url
    
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
        "Accept-Encoding": "gzip, deflate, br",
        "Connection": "keep-alive",
        "Referer": "https://orzice.com/",
        "Sec-Ch-Ua": "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"",
        "Sec-Ch-Ua-Mobile": "?0",
        "Sec-Ch-Ua-Platform": "\"Windows\"",
        "Sec-Fetch-Dest": "document",
        "Sec-Fetch-Mode": "navigate",
        "Sec-Fetch-Site": "same-origin",
        "Sec-Fetch-User": "?1",
        "Upgrade-Insecure-Requests": "1"
    }
    
    try:
        session = requests.Session()
        session.headers.update(headers)
        
        print("正在发送请求...")
        response = session.get(url, timeout=30)
        response.raise_for_status()
        print(f"响应状态码: {response.status_code}")
        print(f"响应内容长度: {len(response.text)}")
        
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, 'html.parser')
            gun_data = []
            
            # 从表格中提取数据
            table_rows = soup.find_all('tr')
            for row in table_rows:
                item_data = {}
                
                # 提取名称
                name_elem = row.find(['td', 'div', 'span'], class_=re.compile(r'(name|title|item-name|weapon-name)'))
                if not name_elem:
                    name_elem = row.find(['td', 'div'], string=re.compile(r'(步枪|冲锋枪|轻机枪|狙击枪|手枪|霰弹枪)'))
                
                if name_elem:
                    item_data['name'] = name_elem.get_text(strip=True)
                
                # 提取价格 - 匹配 price-cell 类下的 icon-gold
                price_cell = row.find('td', class_='price-cell')
                if price_cell:
                    price_elem = price_cell.find('span', class_='icon-gold')
                    if price_elem:
                        price_text = price_elem.get_text(strip=True)
                        # 清理模板语法，提取数字
                        price_match = re.search(r"(\d{1,3}(?:,\d{3})*)", price_text)
                        if price_match:
                            item_data['price'] = price_match.group(1)
                        else:
                            item_data['price'] = price_text
                
                # 如果没找到，尝试其他价格选择器
                if 'price' not in item_data:
                    price_elem = row.find(['span', 'div'], class_=re.compile(r'(icon-gold|price|cost|num)'))
                    if price_elem:
                        price_text = price_elem.get_text(strip=True)
                        price_match = re.search(r"(\d{1,3}(?:,\d{3})*)", price_text)
                        if price_match:
                            item_data['price'] = price_match.group(1)
                        else:
                            item_data['price'] = price_text
                
                # 提取品质/等级
                grade_elem = row.find(['span', 'div', 'td'], class_=re.compile(r'(grade|level|quality|品质)'))
                if grade_elem:
                    item_data['grade'] = grade_elem.get_text(strip=True)
                
                # 提取图标
                icon_elem = row.find('img')
                if icon_elem:
                    item_data['icon'] = icon_elem.get('src', '')
                
                if item_data.get('name') or item_data.get('price'):
                    gun_data.append(item_data)
            
            # 从其他元素提取数据作为补充
            items = soup.find_all(['div', 'li'], class_=re.compile(r'(item|gun|weapon|list)'))
            for item in items:
                item_data = {}
                
                name_elem = item.find(['span', 'h3', 'h4', 'div'], class_=re.compile(r'(name|title)'))
                if name_elem:
                    item_data['name'] = name_elem.get_text(strip=True)
                
                # 查找价格
                price_cell = item.find('td', class_='price-cell')
                if price_cell:
                    price_elem = price_cell.find('span', class_='icon-gold')
                    if price_elem:
                        price_text = price_elem.get_text(strip=True)
                        price_match = re.search(r"(\d{1,3}(?:,\d{3})*)", price_text)
                        if price_match:
                            item_data['price'] = price_match.group(1)
                        else:
                            item_data['price'] = price_text
                
                if 'price' not in item_data:
                    price_elem = item.find(['span', 'div'], class_=re.compile(r'(icon-gold|price|cost|price.*num)'))
                    if price_elem:
                        price_text = price_elem.get_text(strip=True)
                        price_match = re.search(r"(\d{1,3}(?:,\d{3})*)", price_text)
                        if price_match:
                            item_data['price'] = price_match.group(1)
                        else:
                            item_data['price'] = price_text
                
                grade_elem = item.find(['span', 'div'], class_=re.compile(r'(grade|level|quality)'))
                if grade_elem:
                    item_data['grade'] = grade_elem.get_text(strip=True)
                
                icon_elem = item.find('img')
                if icon_elem:
                    item_data['icon'] = icon_elem.get('src', '')
                
                if item_data.get('name') or item_data.get('price'):
                    gun_data.append(item_data)
            
            # 过滤掉没有价格的数据并去重
            unique_data = []
            seen = set()
            for item in gun_data:
                if item.get('name'):
                    # 优先保留有价格的数据
                    key = str(item.get('name', ''))
                    if key not in seen:
                        seen.add(key)
                        unique_data.append(item)
                    elif 'price' in item:
                        # 如果已存在但当前项有价格，替换掉
                        idx = next(i for i, d in enumerate(unique_data) if d.get('name') == key)
                        if 'price' not in unique_data[idx]:
                            unique_data[idx] = item
            
            return unique_data
        
        else:
            print(f"请求失败，状态码: {response.status_code}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"请求异常: {e}")
        return None

# def save_to_file(data, filename='gun_data.json'):
#     if data:
#         with open(filename, 'w', encoding='utf-8') as f:
#             json.dump(data, f, ensure_ascii=False, indent=2)
#         print(f"数据已保存到 {filename}")

def get_firearms():
    print("=== 开始抓取枪械数据 ===")
    page = 1
    while True:
        print(f"当前页为第{page}页")
        url = f"https://orzice.com/v/zhanbei?a=zhanbei&top=2-1&p={page}&grade=-1&n=%E6%9E%AA%E6%A2%B0"
        gun_data = get_gun_data(url)
        page += 1
        if gun_data:
            print(f"\n成功抓取到 {len(gun_data)} 条枪械数据")
            print("\n=== 数据预览 ===")
            for i, item in enumerate(gun_data):
                print(f"\n{i+1}. {item.get('name', '未知')},  价格: {item['price']}")
            # save_to_file(gun_data)
        else:
            print("未能抓取到数据")
            break

def get_accessories():
    print("=== 开始抓取配件数据 ===")
    page = 1
    while True:
        print(f"当前页为第{page}页")
        url = f"https://orzice.com/v/zhanbei?a=zhanbei&top=2-1&p={page}&grade=-1&n=%E9%85%8D%E4%BB%B6"
        gun_data = get_gun_data(url)
        page += 1
        if gun_data:
            print(f"\n成功抓取到 {len(gun_data)} 条配件数据")
            print("\n=== 数据预览 ===")
            for i, item in enumerate(gun_data):
                print(f"\n{i+1}. {item.get('name', '未知')},  价格: {item['price']}")
            # save_to_file(gun_data)
        else:
            print("未能抓取到数据")
            break


def main():
    get_firearms()
    get_accessories()
if __name__ == "__main__":
    main()