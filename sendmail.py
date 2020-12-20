import re
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.header import Header
from email.mime.application import MIMEApplication

def send_mail(text,fileUrl):
    # 发信方的信息：发信邮箱，QQ 邮箱授权码
    from_addr = '601337784@qq.com'
    password = 'sfiqnsvrwtlqbdgj'
    # 收信方邮箱
    to_addr = 'sureqweasdzxc_367@163.com'
    # 发信服务器
    smtp_server = 'smtp.qq.com'

    msg = MIMEText(text, 'plain', 'utf-8')
    x = re.search('/(.*[csv,pdf])',fileUrl)   # 通过路径解析文件名
    print(x.group()[1:])
    fileName = x.group()[1:]
    pdfApart = MIMEApplication(open(fileUrl, 'rb').read())
    pdfApart.add_header('Content-Disposition', 'attachment', filename=fileName)

    msg_file = MIMEMultipart()
    msg_file['From'] = Header(from_addr)
    msg_file['To'] = Header(to_addr)

    msg_file.attach(msg)
    msg_file.attach(pdfApart)
    msg_file['Subject'] = 'title1'

    # 开启发信服务，这里使用的是加密传输
    server = smtplib.SMTP_SSL(smtp_server)
    server.connect(smtp_server, 465)
    # 登录发信邮箱
    server.login(from_addr, password)
    # 发送邮件
    server.sendmail(from_addr, to_addr, msg_file.as_string())
    # 关闭服务器
    server.quit()


if __name__=="__main__":
    send_mail(
        '''
            申请邮件
            ObjectId:"5fdf354958a7eb58164ca245",
            
            添加pdf邮件
        ''',
        "E:/账号如何开通.pdf"
    )