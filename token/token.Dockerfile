FROM centos:7

RUN cp ./refreshToken.py /app && \
    cp ./chromedriver /app

WORKDIR /app

#ChromeDriver 103.0.5060.134
RUN yum -y update && yum -y install libX11-1.6.7-4.el7_9.x86_64 --setopt=protected_multilib=false && \
    yum -y install GConf2 && yum -y install fontconfig-devel && \
    yum -y install epel-release && yum -y install chromium && \
    yum list installed | grep chro && chmod +x chromedriver && \
    ln -s /app/chromedriver /usr/local/bin/chromedriver && \
    ln -s /app/chromedriver /usr/bin/chromedriver && \
    yum -y install python3-pip && \
    pip3 install selenium

EXPOSE 80

CMD python3 /app/refreshToken.py