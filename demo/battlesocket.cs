using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using BestHTTP.SocketIO;

public class BattleSocket
{
    private SocketManager manager = null;
    private Socket socket = null;

    private BattleManager battleManager = null;
    private BattleManager BattleManager {
        get {
            return battleManager;
        }
        set {
            battleManager = value;
        }
    }

    private bool connected = false;
    private int lastSequence = -1;
    private float totalElapsedTime = 0.0f;

    private bool started = false;
    public bool Started {
        get {
            return started;
        }
        private set {
            started = value;
        }
    }

    //private List<>

    public  BattleSocket()
    {
        BattleManager = BattleManager.Instance;

        SocketOptions options = new SocketOptions();
        options.AutoConnect = false;
        options.ConnectWith = BestHTTP.SocketIO.Transports.TransportTypes.WebSocket;
        manager = new SocketManager(new System.Uri(DataManager.Instance.socketHost), options);
        socket = manager.GetSocket("/battle");

        socket.On(SocketIOEventTypes.Connect, (Socket _socket, Packet packet, object[] args) => {
            Dictionary<string, string> parameters = new Dictionary<string, string>();
            parameters.Add("battleId", BattleManager.BattleProperties.ID);
            socket.Emit("battle", parameters);
        });

        socket.On("battle", (Socket _socket, Packet packet, object[] args) => {
            Debug.Log("connected");

            connected = true;
        });

        socket.On(SocketIOEventTypes.Disconnect, (Socket _socket, Packet packet, object[] args) => {
            connected = false;
        });

        socket.On(SocketIOEventTypes.Error, (Socket _socket, Packet packet, object[] args) => {
            Error error = args[0] as Error;
            Debug.Log(error.ToString());
        });

        socket.On("battleEvent", (Socket _socket, Packet packet, object[] args) => {
            Dictionary<string, object> data = args[0] as Dictionary<string, object>;
            string key = data["key"] as string;
            int sequence = System.Int32.Parse(data["sequence"] as string);
            int delay = System.Int32.Parse(data["delay"] as string);

            BattleManager.StartCoroutine(processBattleEvent(key, sequence, delay, data));
        });

        manager.Open();
    }

    private IEnumerator processBattleEvent(string key, int sequence, int delay, Dictionary<string, object> data) {
        if (sequence <= lastSequence) {
            yield break;
        }

        while (sequence != lastSequence + 1) {
            yield return 0;
        }

        float elapsed = 0.0f;

        while (elapsed < delay) {
            elapsed += Time.deltaTime * 1000;

            yield return 0;
        }

        Debug.Log(string.Format("processing battle event - {0} - {1} {2}", sequence, key, totalElapsedTime));

        PartyController opponentPartyControlller = BattleManager.getOpponentPartyController();
        BattleController battleController = BattleManager.getCurrentBattleController();

        if (key.Equals("deal_cards")) {
            List<string> dealedCardIds = new List<string>();
            (data["dealed_card_ids"] as List<object>).ForEach((cardId) =>
            {
                dealedCardIds.Add(cardId as string);
            });

            List<BaseCard> cards = battleController.grabCards(dealedCardIds);

            BattleManager.sendEvent("dealCards", cards);
        } else if (key.Equals("use_card")) {
            string cardId = data["card_id"].ToString();

            List<string> targetIds = new List<string>();
            (data["target_ids"] as List<object>).ForEach((targetId) =>
            {
                targetIds.Add(targetId as string);
            });
            List<BaseActor> targets = opponentPartyControlller.grabActors(targetIds);

            BaseCard card = battleController.grabCard(cardId);

            BattleManager.sendEvent("useCard", card, battleController.Actor, targets);
        } else if (key.Equals("end_turn")) {
            BattleManager.sendEvent("startBattleTurn");
        } else if (key.Equals("replenish_cards")) {
            BattleManager.sendEvent("replenishCards");
        } else if (key.Equals("next_party_stage")) {
            BattleManager.sendEvent("nextPartyStage");
        } else if (key.Equals("end")) {
            List<NodeRewardProperties> nodeRewards = new List<NodeRewardProperties>();
            (data["node_rewards"] as List<object>).ForEach((nodeReward) =>
            {
                nodeRewards.Add(new NodeRewardProperties(nodeReward as Dictionary<string, object>));
            });

            Dictionary<string, int> rewardRows = new Dictionary<string, int>();
            foreach (KeyValuePair<string, object> kvp in data["reward_rows"] as Dictionary<string, object>)
            {
                rewardRows.Add(kvp.Key, System.Convert.ToInt32(kvp.Value));
            };

            BattleManager.sendEvent("end", nodeRewards, rewardRows);
        }

        lastSequence = sequence;
    }

    public void close()
    {
        manager.Close();
    }

    private IEnumerator invokeNextFrame(System.Action action) {
        yield return 0;

        Debug.Log(Time.realtimeSinceStartup);

        action();
    }

    private IEnumerator invokeNextFrame<T>(System.Action<T> action, T arg1) {
        yield return 0;

        Debug.Log(Time.realtimeSinceStartup);

        action(arg1);
    }

    public void dispatchStart()
    {
        if (connected == false) {
            BattleManager.StartCoroutine(invokeNextFrame(dispatchStart));

            return;
        }

        string url = string.Format("/battles/{0}/start.json", BattleManager.BattleProperties.ID);

        NetworkManager.Instance.post<BattleProperties.BattleWrapper>(url, null).Then((wrapper) => {
            Started = true;
        }).Catch((exception) => {
            GameManager.Instance.showPopup("UI/ErrorPopup");
        });
    }

    public void dispatchUseCard(BattleCardController cardController)
    {
        if (connected == false) {
            BattleManager.StartCoroutine(invokeNextFrame<BattleCardController>(dispatchUseCard, cardController));

            return;
        }

        string url = string.Format("/battles/{0}/use_card.json", BattleManager.BattleProperties.ID);

        foreach(BattleCardController battleCardController in BattleManager.CardsDeckController.DealedCardControllers) {
            battleCardController.Button.enabled = false;
        }

        BattleProperties battleProperties = new BattleProperties
        {
            CardID = cardController.Card.Properties.ID,
        };
        NetworkManager.Instance.post<BattleProperties.BattleWrapper>(url, battleProperties.toJson()).Then((wrapper) => {
            foreach (BattleCardController battleCardController in BattleManager.CardsDeckController.DealedCardControllers)
            {
                if (cardController != battleCardController)
                {
                    battleCardController.Button.enabled = true;
                }
            }
        }).Catch((exception) => {
            foreach (BattleCardController battleCardController in BattleManager.CardsDeckController.DealedCardControllers)
            {
                battleCardController.Button.enabled = true;
            }

            GameManager.Instance.showPopup("UI/ErrorPopup");
        });
    }

    public void dispatchEndTurn(Button button)
    {
        string url = string.Format("/battles/{0}/end_turn.json", BattleManager.BattleProperties.ID);

        button.enabled = false;

        NetworkManager.Instance.post<BattleProperties.BattleWrapper>(url, null).Then((wrapper) => {
            button.enabled = true;
        }).Catch((exception) => {
            button.enabled = true;

            GameManager.Instance.showPopup("UI/ErrorPopup");
        });
    }

    public void update() {
        totalElapsedTime += Time.deltaTime;
    }
}
