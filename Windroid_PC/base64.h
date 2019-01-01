#ifndef BASE64_H
#define BASE64_H

#include <QString>

class Base64
{
public:
    /*
     * 功能：静态成员函数，将字节数组转换为Base64编码字符串
     * 参数说明：
     *      binaryData：要转换的字节数组
     * 返回值：
     *      转换后得到的Base64编码字符串
     * 异常抛出：
     *      无
     * 说明：
     *      空说明列表指出函数不抛出任何异常
     *      如果一个函数声明没有指定异常说明，则该函数可以抛出任意类型的异常
     */
    static QString encode(const QByteArray & binaryData) throw();

    /*
     * 功能：静态成员函数，将Base64编码字符串解码为字节数组
     * 参数说明：
     *      base64String：要转换的Base64编码字符串
     * 返回值：
     *      解码后得到的字节数组
     * 异常抛出：
     *      抛出整型异常
     *          -1：数据错误
     * 说明：
     *      字符串中允许任意的空白字符、回车换行符、连字符
     */
    static QByteArray decode(const QString & base64String) throw(int);
};

#endif // BASE64_H
