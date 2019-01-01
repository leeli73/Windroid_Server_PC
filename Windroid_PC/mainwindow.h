#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "qclipboard"
#include "QString"
#include "QMessageBox"
#include "QObject"
#include "QUrl"
#include "QByteArray"
#include "QVariant"
#include "QDebug"
#include "QtNetwork/QNetworkReply"
#include "QtNetwork/QNetworkRequest"
#include "QtNetwork/QNetworkAccessManager"
#include "base64.h"

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    QString UserID;
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();
    QNetworkAccessManager *m_accessManager;
    bool Login();
    bool ConnectServer();
    QString GetClipboardText();
    void SetClipboardText(QString text);
    void GetData();
    void SetData();


private slots:
    void on_pushButton_clicked();
    void finishedSlot(QNetworkReply *reply);

private:
    Ui::MainWindow *ui;
};

#endif // MAINWINDOW_H
