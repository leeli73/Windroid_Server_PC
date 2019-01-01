#include "mainwindow.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    setMinimumSize(380, 251);
    setMaximumSize(380, 251);
    m_accessManager = new QNetworkAccessManager(this);
    QObject::connect(m_accessManager, SIGNAL(finished(QNetworkReply*)), this, SLOT(finishedSlot(QNetworkReply*)));
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_pushButton_clicked()
{
    //QMessageBox::warning(NULL, "警告", GetClipboardText());
    Login();
}
QString MainWindow::GetClipboardText()
{
    QClipboard *clipboard = QApplication::clipboard();
    return clipboard->text();
}
void MainWindow::SetClipboardText(QString text)
{
    QClipboard *clipboard = QApplication::clipboard();
    clipboard->setText(text);
}
bool MainWindow::Login()
{
    QString Username = ui->lineEdit_Username->text();
    QString Password = ui->lineEdit_Password->text();
    if(Username.isEmpty()||Password.isEmpty())
    {
        QMessageBox::warning(NULL, "警告", "账号和密码不能为空！");
        return false;
    }
    Base64 base64;
    Username = base64.encode(Username.toLatin1());
    Password = base64.encode(Password.toLatin1());
    QNetworkRequest request;
    request.setUrl(QUrl("http://127.0.0.1:6888/Login"));
    request.setHeader(QNetworkRequest::ContentTypeHeader,"application/x-www-form-urlencoded");
    request.setRawHeader("Accept","text/html, application/xhtml+xml, */*");
    request.setRawHeader("Accept-Language","zh-CN");
    request.setRawHeader("X-Requested-With","XMLHttpRequest");
    request.setRawHeader("User-Agent","Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)");
    request.setRawHeader("Content-Type","application/x-www-form-urlencoded");
    request.setRawHeader("Accept-Encoding","gzip,deflate");
    request.setRawHeader("Connection","Keep-Alive");
    request.setRawHeader("Cache-Control","no-cache");
    QByteArray postData;
    postData.append("Username="+Username+"&Password="+Password);
    QNetworkReply* reply = m_accessManager->post(request,postData);
    return true;
}
bool MainWindow::ConnectServer()
{
    return true;
}
void MainWindow::finishedSlot(QNetworkReply *reply)
{
    if (reply->error() == QNetworkReply::NoError)
    {
        QByteArray bytes = reply->readAll();
        qDebug()<<bytes;
        QString string = QString::fromUtf8(bytes);
        QStringList temp = string.split("@");
        if(temp[0] == "LoginSuccess")
        {
            QMessageBox::about(NULL, "About", "登录成功");
            UserID = temp[1];
            this->hide();
            while(1)
            {
                GetData();
                SetData();
            }
        }
        else if(temp[0] == "New")
        {
            SetClipboardText(temp[1]);
        }
    }
    else
    {
        qDebug()<<"handle errors here";
        QVariant statusCodeV = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute);
        //statusCodeV是HTTP服务器的相应码，reply->error()是Qt定义的错误码，可以参考QT的文档
        qDebug( "found error ....code: %d %d\n", statusCodeV.toInt(), (int)reply->error());
        qDebug(qPrintable(reply->errorString()));
    }
    reply->deleteLater();

}
void MainWindow::GetData()
{
    QString myUserID = UserID;
    Base64 base64;
    myUserID = base64.encode(myUserID.toLatin1());
    QNetworkRequest request;
    request.setUrl(QUrl("http://127.0.0.1:6888/GetData"));
    request.setHeader(QNetworkRequest::ContentTypeHeader,"application/x-www-form-urlencoded");
    request.setRawHeader("Accept","text/html, application/xhtml+xml, */*");
    request.setRawHeader("Accept-Language","zh-CN");
    request.setRawHeader("X-Requested-With","XMLHttpRequest");
    request.setRawHeader("User-Agent","Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)");
    request.setRawHeader("Content-Type","application/x-www-form-urlencoded");
    request.setRawHeader("Accept-Encoding","gzip,deflate");
    request.setRawHeader("Connection","Keep-Alive");
    request.setRawHeader("Cache-Control","no-cache");
    QByteArray postData;
    postData.append("UserID="+myUserID);
    QNetworkReply* reply = m_accessManager->post(request,postData);
}
void MainWindow::SetData()
{
    QString myUserID = UserID;
    Base64 base64;
    myUserID = base64.encode(myUserID.toLatin1());
    QNetworkRequest request;
    request.setUrl(QUrl("http://127.0.0.1:6888/SetData"));
    request.setHeader(QNetworkRequest::ContentTypeHeader,"application/x-www-form-urlencoded");
    request.setRawHeader("Accept","text/html, application/xhtml+xml, */*");
    request.setRawHeader("Accept-Language","zh-CN");
    request.setRawHeader("X-Requested-With","XMLHttpRequest");
    request.setRawHeader("User-Agent","Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)");
    request.setRawHeader("Content-Type","application/x-www-form-urlencoded");
    request.setRawHeader("Accept-Encoding","gzip,deflate");
    request.setRawHeader("Connection","Keep-Alive");
    request.setRawHeader("Cache-Control","no-cache");
    QByteArray postData;
    postData.append("UserID="+myUserID+"Data="+base64.encode(GetClipboardText().toLatin1()));
    QNetworkReply* reply = m_accessManager->post(request,postData);
}
