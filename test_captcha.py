#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import json
import base64
from io import BytesIO
from PIL import Image

# API基础URL
BASE_URL = "http://localhost:8080/api/v1"

def test_captcha_flow():
    """测试完整的验证码流程"""
    print("=== 测试验证码功能 ===")
    
    # 1. 获取验证码
    print("\n1. 获取验证码...")
    response = requests.get(f"{BASE_URL}/auth/captcha")
    if response.status_code == 200:
        captcha_data = response.json()
        print(f"✅ 验证码获取成功")
        print(f"   验证码ID: {captcha_data['captcha_id']}")
        print(f"   认证Token: {captcha_data['auth_token']}")
        
        # 显示验证码图片
        captcha_b64 = captcha_data['captcha_b64']
        captcha_id = captcha_data['captcha_id']
        auth_token = captcha_data['auth_token']
        
        # 解码并保存验证码图片
        try:
            # 检查是否有data:image前缀
            if captcha_b64.startswith('data:image'):
                # 移除前缀，只保留base64部分
                captcha_b64 = captcha_b64.split(',', 1)[1]
            
            image_data = base64.b64decode(captcha_b64)
            with open("captcha.png", "wb") as f:
                f.write(image_data)
            print(f"   验证码图片已保存为 captcha.png")
            
            # 显示图片信息
            try:
                from PIL import Image
                img = Image.open("captcha.png")
                print(f"   图片大小: {img.size}, 格式: {img.format}")
            except:
                pass
        except Exception as e:
            print(f"   保存验证码图片失败: {e}")
        
        return captcha_id, auth_token
    else:
        print(f"❌ 获取验证码失败: {response.status_code}")
        print(f"   错误信息: {response.text}")
        return None, None

def test_verify_captcha(captcha_id, auth_token):
    """测试验证码验证"""
    print("\n2. 验证验证码...")
    print("   请输入验证码 (查看captcha.png文件):")
    captcha_code = input("   验证码: ").strip()
    
    data = {
        "captcha_id": captcha_id,
        "captcha_code": captcha_code
    }
    
    response = requests.post(f"{BASE_URL}/auth/verify-captcha", json=data)
    if response.status_code == 200:
        result = response.json()
        print(f"✅ 验证码验证成功")
        print(f"   新的认证Token: {result['auth_token']}")
        return result['auth_token']
    else:
        print(f"❌ 验证码验证失败: {response.status_code}")
        print(f"   错误信息: {response.text}")
        return None

def test_login_with_auth_token(auth_token):
    """测试使用认证token登录"""
    print("\n3. 测试登录 (需要认证token)...")
    
    # 这里需要真实的用户名和密码
    print("   请输入测试用户名:")
    username = input("   用户名: ").strip()
    print("   请输入测试密码:")
    password = input("   密码: ").strip()
    
    data = {
        "username": username,
        "password": password,
        "auth_token": auth_token
    }
    
    response = requests.post(f"{BASE_URL}/auth/login", json=data)
    if response.status_code == 200:
        result = response.json()
        print(f"✅ 登录成功")
        print(f"   JWT Token: {result['token']}")
        print(f"   用户信息: {result['user']}")
    else:
        print(f"❌ 登录失败: {response.status_code}")
        print(f"   错误信息: {response.text}")

def test_register_with_auth_token(auth_token):
    """测试使用认证token注册"""
    print("\n4. 测试注册 (需要认证token)...")
    
    print("   请输入注册信息:")
    username = input("   用户名: ").strip()
    email = input("   邮箱: ").strip()
    password = input("   密码: ").strip()
    invite_code = input("   邀请码: ").strip()
    
    data = {
        "username": username,
        "email": email,
        "password": password,
        "invite_code": invite_code,
        "auth_token": auth_token
    }
    
    response = requests.post(f"{BASE_URL}/auth/register", json=data)
    if response.status_code == 201:
        result = response.json()
        print(f"✅ 注册成功")
        print(f"   JWT Token: {result['token']}")
        print(f"   用户信息: {result['user']}")
    else:
        print(f"❌ 注册失败: {response.status_code}")
        print(f"   错误信息: {response.text}")

def main():
    """主测试函数"""
    try:
        # 测试验证码获取
        captcha_id, auth_token = test_captcha_flow()
        if not captcha_id or not auth_token:
            print("❌ 验证码获取失败，测试终止")
            return
        
        # 测试验证码验证
        new_auth_token = test_verify_captcha(captcha_id, auth_token)
        if not new_auth_token:
            print("❌ 验证码验证失败，测试终止")
            return
        
        # 测试登录
        test_login_with_auth_token(new_auth_token)
        
        # 测试注册
        test_register_with_auth_token(new_auth_token)
        
        print("\n=== 测试完成 ===")
        
    except KeyboardInterrupt:
        print("\n\n测试被用户中断")
    except Exception as e:
        print(f"\n测试过程中发生错误: {e}")

if __name__ == "__main__":
    main() 