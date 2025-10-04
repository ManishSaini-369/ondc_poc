import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '10s',
};

const subscriberData = [
    "Amazon.com", "Paytm.com", "Oyo.com", "Usha.com", "Zomato.com", "Swiggy.com", "Flipkart.com", "Myntra.com", "Meesho.com", "TataNeu.com",
    "RelianceRetail.com", "JioMart.com", "BigBasket.com", "Blinkit.com", "Spencers.com", "DMart.com", "Nykaa.com", "Ajio.com", "Snapdeal.com", "ShopClues.com",
    "Pepperfry.com", "FirstCry.com", "UrbanLadder.com", "Lenskart.com", "PharmEasy.com", "1mg.com", "ApolloPharmacy.com", "MedPlus.com", "NetMeds.com", "RedBus.com",
    "MakeMyTrip.com", "Yatra.com", "EaseMyTrip.com", "Cleartrip.com", "Ixigo.com", "BookMyShow.com", "Gaana.com", "Spotify.com", "Wynk.com", "Hotstar.com",
    "SonyLiv.com", "Zee5.com", "Voot.com", "ALTBalaji.com", "MXPlayer.com", "JioSaavn.com", "Cars24.com", "CarDekho.com", "Spinny.com", "Ola.com",
    "Uber.com", "Rapido.com", "Bounce.com", "Zoomcar.com", "Porter.com", "Delhivery.com", "EcomExpress.com", "Shadowfax.com", "BlueDart.com", "DTDC.com",
    "XpressBees.com", "Shiprocket.com", "FedEx.com", "DHL.com", "Maersk.com", "AirIndia.com", "IndiGo.com", "SpiceJet.com", "Vistara.com", "AkasaAir.com",
    "GoFirst.com", "IRCTC.com", "HDFC.com", "ICICI.com", "SBI.com", "AxisBank.com", "Kotak.com", "YesBank.com", "IDFCFirst.com", "PNB.com",
    "BOB.com", "HSBC.com", "CitiBank.com", "StandardChartered.com", "Razorpay.com", "PhonePe.com", "GooglePay.com", "BharatPe.com", "Freecharge.com", "Mobikwik.com",
    "PayU.com", "PolicyBazaar.com", "LIC.com", "HDFCLife.com", "ICICILombard.com", "StarHealth.com", "MaxBupa.com", "OrientalInsurance.com", "Boat.com", "Noise.com"
];

const url = 'http://localhost:8081/lookup';
const params = {
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer AAAAAAAAAAAAAAAAAAAAALnsihDAB+rg4iHoXcf+abZL3l0T6ajj2x78Zj0nDA3vY4tvFmlCp+qCOJCDQvHJfWtJm21/FdBdRRVlrDzou1MopkiYtKh0bAllmEdBxYecuLZZDin1LJ5pzoEJXqqWwEoD689MfZzGCvt3lt/Dvajt7o1TFeLjCEvRE8T/iTifG84M',
  },
};

export default function () {
  // Select a random subscriber from the data array
  const subscriberId = subscriberData[Math.floor(Math.random() * subscriberData.length)];

  const payload = JSON.stringify({
    subscriber_id: subscriberId,
  });

  const res = http.post(url, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1); // Wait for 1 second between requests
}
