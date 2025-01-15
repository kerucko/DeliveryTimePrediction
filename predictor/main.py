from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
import joblib
from kafka import KafkaConsumer, KafkaProducer
import json
import threading
import logging


app = FastAPI()


class DeliveryRequest(BaseModel):
    Order_ID: str
    Distance_km: float
    Weather: str
    Traffic_Level: str
    Time_of_Day: str
    Vehicle_Type: str
    Preparation_Time_min: int
    Courier_Experience_yrs: float


MODEL_PATH = "delivery_time_model.pkl"
KAFKA_BROKER = "kafka:29092"
INPUT_TOPIC = "tasks"
OUTPUT_TOPIC = "completed"


def load_model():
    try:
        model = joblib.load(MODEL_PATH)
        logging.info(f"Model loaded from {MODEL_PATH}")
        return model
    except FileNotFoundError:
        raise RuntimeError(f"Model file not found at {MODEL_PATH}")

# Initialize Kafka consumer and producer
def initialize_kafka():
    consumer = KafkaConsumer(
        INPUT_TOPIC,
        bootstrap_servers=KAFKA_BROKER,
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )

    producer = KafkaProducer(
        bootstrap_servers=KAFKA_BROKER,
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )

    return consumer, producer


def process_kafka_messages(consumer, producer, model):
    for message in consumer:
        try:
            data = message.value
            delivery_request = DeliveryRequest(**data)

            input_data = pd.DataFrame([{
                'Distance_km': delivery_request.Distance_km,
                'Weather': delivery_request.Weather,
                'Traffic_Level': delivery_request.Traffic_Level,
                'Time_of_Day': delivery_request.Time_of_Day,
                'Vehicle_Type': delivery_request.Vehicle_Type,
                'Preparation_Time_min': delivery_request.Preparation_Time_min,
                'Courier_Experience_yrs': delivery_request.Courier_Experience_yrs
            }])

            prediction = model.predict(input_data)[0]
            result = {
                "Order_ID": delivery_request.Order_ID, 
                "predicted_delivery_time": float(prediction)
            }


            producer.send(OUTPUT_TOPIC, result)
            logging.info(f"Delivery time for order {delivery_request.Order_ID} sent to Kafka.")
        
        except Exception as e:
            logging.error(f"Error processing message: {e}")


def start_kafka_listener():
    consumer, producer = initialize_kafka()
    model = load_model()
    kafka_thread = threading.Thread(target=process_kafka_messages, args=(consumer, producer, model), daemon=True)
    kafka_thread.start()


@app.get("/")
async def root():
    return {"message": "Delivery Time Prediction Service is running."}


start_kafka_listener()
