FROM ubuntu:23.10

WORKDIR /app

RUN apt-get update && apt-get install -y \
    python3-venv \
    python3-pip \
    golang \
    libgirepository1.0-dev \
    imagemagick \
    libmagick++-dev \
    ffmpeg

RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

COPY requirements.txt .
RUN pip install --upgrade wheel
RUN pip install -r requirements.txt

COPY libmagick-config.txt .
RUN cat libmagick-config.txt > /etc/ImageMagick-6/policy.xml

COPY build/binapp /app/bin/binapp

COPY .env .
COPY pkg/ pkg/
COPY resources/ resources/
COPY generated/ generated/

ADD ./generated /app/generated

EXPOSE 1323

CMD ["bin/binapp"]
