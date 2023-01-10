import requests


def get_tradeable_pairs(pairs):
    tradeable_pairs = []

    for pair in pairs:
        # add pair name to object
        pairs[pair]["pair_name"] = pair

        # only allows 'online' pairs
        if pairs[pair]['status'] == 'online':
            tradeable_pairs.append(pairs[pair])

    return tradeable_pairs


def structure_arbitrage_pairs(pairs):
    triangular_pairs_list = []
    remove_duplicates_list = []
    pairs_list = pairs[:10]

    # get pair a
    for pair_a in pairs:
        a_base = pair_a['base']
        a_quote = pair_a['quote']
        # assign a to a box
        a_pair_box = [a_base, a_quote]

        # get pair b
        for pair_b in pairs:
            b_base = pair_b['base']
            b_quote = pair_b['quote']
            b_pair_box = [b_base, b_quote]

            # check pair b
            if b_pair_box != a_pair_box:
                if b_base in a_pair_box or b_quote in a_pair_box:

                    # get pair c
                    for pair_c in pairs:
                        c_base = pair_c['base']
                        c_quote = pair_c['quote']
                        c_pair_box = [c_base, c_quote]

                        # count number of matching c items
                        if c_pair_box != a_pair_box and c_pair_box != b_pair_box:
                            combine_all = [pair_a, pair_b, pair_c]
                            pair_box = [a_base, a_quote, b_base,
                                        b_quote, c_base, c_quote]

                            counts_c_base = 0
                            counts_c_quote = 0
                            for i in pair_box:
                                if i == c_base:
                                    counts_c_base += 1
                                if i == c_quote:
                                    counts_c_quote += 1

                            # determine triangular match
                            if counts_c_base == 2 and counts_c_quote == 2 and c_base != c_quote:
                                # combined = sorted(a_pair_box), sorted(
                                #     b_pair_box), sorted(c_pair_box)

                                combined = [a_base + '/' + a_quote, b_base +
                                            "/" + b_quote, c_base + "/" + c_quote]
                                unique_item = sorted(combined)
                                # print(unique_item)

                                # if unique_item == ['AAVE/ZEUR', 'AAVE/ZUSD', 'ZEUR/ZUSD']:
                                #     print('match!')

                                if unique_item not in remove_duplicates_list:
                                    match_dict = {
                                        "a_base": a_base,
                                        "b_base": b_base,
                                        "c_base": c_base,
                                        "a_quote": a_quote,
                                        "b_quote": b_quote,
                                        "c_quote": c_quote,
                                        "a_pair": a_pair_box,
                                        "b_pair": b_pair_box,
                                        "c_pair": c_pair_box,
                                        "a_altname": pair_a['altname'],
                                        "b_altname": pair_b['altname'],
                                        "c_altname": pair_c['altname'],
                                        "unique_item": unique_item,
                                        "a_pair_name": pair_a['pair_name'],
                                        "b_pair_name": pair_b['pair_name'],
                                        "c_pair_name": pair_c['pair_name'],
                                    }

                                    triangular_pairs_list.append(match_dict)
                                    remove_duplicates_list.append(unique_item)
    return triangular_pairs_list


def get_price_for_t_pair(t_pair):
    # get order book for pair
    a_altname = t_pair['a_altname']
    b_altname = t_pair['b_altname']
    c_altname = t_pair['c_altname']

    a_resp = requests.get(
        'https://api.kraken.com/0/public/Depth?pair={}'.format(a_altname)).json()
    b_resp = requests.get(
        'https://api.kraken.com/0/public/Depth?pair={}'.format(b_altname)).json()
    c_resp = requests.get(
        'https://api.kraken.com/0/public/Depth?pair={}'.format(c_altname)).json()

    # print(a_resp)
    # print(b_resp)
    # print(c_resp)

    # extract pairs info
    return {
        "pair_a_ask": a_resp['result'][t_pair['a_pair_name']]['asks'],
        "pair_a_bid": a_resp['result'][t_pair['a_pair_name']]['bids'],
        "pair_b_ask": b_resp['result'][t_pair['b_pair_name']]['asks'],
        "pair_b_bid": b_resp['result'][t_pair['b_pair_name']]['bids'],
        "pair_c_ask": c_resp['result'][t_pair['c_pair_name']]['asks'],
        "pair_c_bid": c_resp['result'][t_pair['c_pair_name']]['bids'],
    }


def calc_tri_arb_surface_rate(t_pair, prices):
    starting_amount = 1
    min_surface_rate = 0
    surface_dict = {}
    contract_2 = ""
    direction_trade_1 = ""
    direction_trade_2 = ""
    direction_trade_3 = ""
    acquired_coin_t2 = 0
    acquired_coin_t3 = 0
    calculated = 0  # if calculated = 1, all 3 trades have been calculated

    # extract pair variables
    a_base = t_pair['a_base']
    a_quote = t_pair['a_quote']
    b_base = t_pair['b_base']
    b_quote = t_pair['b_quote']
    c_base = t_pair['c_base']
    c_quote = t_pair['c_quote']
    a_pair_name = t_pair['a_pair_name']
    b_pair_name = t_pair['b_pair_name']
    c_pair_name = t_pair['c_pair_name']

    # extract price info
    a_ask = prices['pair_a_ask'][0]
    a_bid = prices['pair_a_bid'][0]
    b_ask = prices['pair_b_ask'][0]
    b_bid = prices['pair_b_bid'][0]
    c_ask = prices['pair_c_ask'][0]
    c_bid = prices['pair_c_bid'][0]

    # set directions and loop through
    direction_list = ["forward", "reverse"]
    for direction in direction_list:
        # set variables for swap info
        swap_1 = 0
        swap_2 = 0
        swap_3 = 0
        swap_1_rate = 0
        swap_2_rate = 0
        swap_3_rate = 0

        """
            if swapping coin left (base) -> right (quote) then * (1 / ASK)
            if swapping coin right (quote) -> left (base) then * BID
        """

        # assume starting with a_base and swap for a_quote
        if direction == "forward":
            swap_1 = a_base
            swap_2 = a_quote
            swap_1_rate = 1 / float(a_ask[0])
            direction_trade_1 = "base_to_quote"

        if direction == "reverse":
            swap_1 = a_quote
            swap_2 = a_base
            swap_1_rate = float(a_bid[0])
            direction_trade_1 = "quote_to_base"

        # place first trade
        contract_1 = a_pair_name
        acquired_coin_t1 = starting_amount * swap_1_rate

        """FORWARD"""
        # SCENARIO 1: Check if a_quote (acquired_coin) matches b_quote
        if direction == "forward":
            if a_quote == b_quote and calculated == 0:
                swap_2_rate = float(b_bid[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "quote_to_base"
                contract_2 = b_pair_name

                # if b_base (acquired coin) matches c_base
                if b_base == c_base:
                    swap_3 = c_base
                    swap_3_rate = 1 / float(c_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = c_pair_name

                # if b_base (acquired coin) matches c_quote
                if b_base == c_quote:
                    swap_3 = c_quote
                    swap_3_rate = float(c_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = c_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # SCENARIO 2: Check if a_quote (acquired_coin) matches b_base
        if direction == "forward":
            if a_quote == b_base and calculated == 0:
                swap_2_rate = 1 / float(b_ask[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "base_to_quote"
                contract_2 = b_pair_name

                # if b_quote (acquired coin) matches c_base
                if b_quote == c_base:
                    swap_3 = c_base
                    swap_3_rate = 1 / float(c_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = c_pair_name

                # if b_quote (acquired coin) matches c_quote
                if b_quote == c_quote:
                    swap_3 = c_quote
                    swap_3_rate = float(c_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = c_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # SCENARIO 3: Check if a_quote (acquired_coin) matches c_quote
        if direction == "forward":
            if a_quote == c_quote and calculated == 0:
                swap_2_rate = float(c_bid[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "quote_to_base"
                contract_2 = c_pair_name

                # if c_base (acquired coin) matches b_base
                if c_base == b_base:
                    swap_3 = b_base
                    swap_3_rate = 1 / float(b_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = b_pair_name

                # if c_base (acquired coin) matches b_quote
                if c_base == b_quote:
                    swap_3 = b_quote
                    swap_3_rate = float(b_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = b_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # SCENARIO 4: Check if a_quote (acquired_coin) matches c_base
        if direction == "forward":
            if a_quote == c_base and calculated == 0:
                swap_2_rate = 1 / float(c_ask[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "base_to_quote"
                contract_2 = c_pair_name

                # if c_quote (acquired coin) matches b_base
                if c_quote == b_base:
                    swap_3 = b_base
                    swap_3_rate = 1 / float(b_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = b_pair_name

                # if b_quote (acquired coin) matches b_quote
                if c_quote == b_quote:
                    swap_3 = b_quote
                    swap_3_rate = float(b_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = b_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        """REVERSE"""
        # SCENARIO 5: Check if a_base (acquired_coin) matches b_quote
        if direction == "reverse":
            if a_base == b_quote and calculated == 0:
                swap_2_rate = float(b_bid[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "quote_to_base"
                contract_2 = b_pair_name

                # if b_base (acquired coin) matches c_base
                if b_base == c_base:
                    swap_3 = c_base
                    swap_3_rate = 1 / float(c_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = c_pair_name

                # if b_base (acquired coin) matches c_quote
                if b_base == c_quote:
                    swap_3 = c_quote
                    swap_3_rate = float(c_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = c_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # SCENARIO 6: Check if a_base (acquired_coin) matches b_base
        if direction == "reverse":
            if a_base == b_base and calculated == 0:
                swap_2_rate = 1 / float(b_ask[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "base_to_quote"
                contract_2 = b_pair_name

                # if b_quote (acquired coin) matches c_base
                if b_quote == c_base:
                    swap_3 = c_base
                    swap_3_rate = 1 / float(c_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = c_pair_name

                # if b_quote (acquired coin) matches c_quote
                if b_quote == c_quote:
                    swap_3 = c_quote
                    swap_3_rate = float(c_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = c_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # SCENARIO 7: Check if a_base (acquired_coin) matches c_quote
        if direction == "reverse":
            if a_base == c_quote and calculated == 0:
                swap_2_rate = float(c_bid[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "quote_to_base"
                contract_2 = c_pair_name

                # if c_base (acquired coin) matches b_base
                if c_base == b_base:
                    swap_3 = b_base
                    swap_3_rate = 1 / float(b_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = b_pair_name

                # if c_base (acquired coin) matches b_quote
                if c_base == b_quote:
                    swap_3 = b_quote
                    swap_3_rate = float(b_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = b_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # SCENARIO 8: Check if a_base (acquired_coin) matches c_base
        if direction == "reverse":
            if a_base == c_base and calculated == 0:
                swap_2_rate = 1 / float(c_ask[0])
                acquired_coin_t2 = acquired_coin_t1 * swap_2_rate
                direction_trade_2 = "base_to_quote"
                contract_2 = c_pair_name

                # if c_quote (acquired coin) matches b_base
                if c_quote == b_base:
                    swap_3 = b_base
                    swap_3_rate = 1 / float(b_ask[0])
                    direction_trade_3 = "base_to_quote"
                    contract_3 = b_pair_name

                # if b_quote (acquired coin) matches b_quote
                if c_quote == b_quote:
                    swap_3 = b_quote
                    swap_3_rate = float(b_bid[0])
                    direction_trade_3 = "quote_to_base"
                    contract_3 = b_pair_name

                acquired_coin_t3 = acquired_coin_t2 * swap_3_rate
                calculated = 1

        # if acquired_coin_t3 >= starting_amount:
        #     print(direction, a_pair_name, b_pair_name, c_pair_name,
        #           starting_amount, acquired_coin_t3)

        """ PROFIT LOSS OUTPUT """

        # profit and loss calculations
        profit_loss = acquired_coin_t3 - starting_amount
        profit_loss_perc = (profit_loss / starting_amount) * \
            100 if profit_loss != 0 else 0

        # trade descriptions
        trade_description_1 = f"Start with {swap_1} of {starting_amount}. Swap at {swap_1_rate} for {swap_2} acquiring {acquired_coin_t1}"
        trade_description_2 = f"Swap {acquired_coin_t1} of {swap_2} at {swap_2_rate} for {swap_3} acquiring {acquired_coin_t2}"
        trade_description_3 = f"Swap {acquired_coin_t2} of {swap_3} at {swap_3_rate} for {swap_1} acquiring {acquired_coin_t3}"

        # if profit_loss > 0:
        #     print("NEW TRADE")
        #     print(trade_description_1)
        #     print(trade_description_2)
        #     print(trade_description_3)

        # output results
        if profit_loss_perc > min_surface_rate:
            surface_dict = {
                "swap_1": swap_1,
                "swap_2": swap_2,
                "swap_3": swap_3,
                "contract_1": contract_1,
                "contract_2": contract_2,
                "contract_3": contract_3,
                "direction_trade_1": direction_trade_1,
                "direction_trade_2": direction_trade_2,
                "direction_trade_3": direction_trade_3,
                "starting_amount": starting_amount,
                "acquired_coin_t1": acquired_coin_t1,
                "acquired_coin_t2": acquired_coin_t2,
                "acquired_coin_t3": acquired_coin_t3,
                "swap_1_rate": swap_1_rate,
                "swap_2_rate": swap_2_rate,
                "swap_3_rate": swap_3_rate,
                "profit_loss": profit_loss,
                "profit_loss_perc": profit_loss_perc,
                "direction": direction,
                "trade_description_1": trade_description_1,
                "trade_description_2": trade_description_2,
                "trade_description_3": trade_description_3
            }

            return surface_dict

    return surface_dict


def reformatted_orderbook(asks, bids, c_direction):
    price_list_main = []

    if c_direction == "base_to_quote":
        for p in asks:
            ask_price = float(p[0])
            adj_price = 1 / ask_price if ask_price != 0 else 0
            adj_quantity = float(p[1]) * ask_price
            price_list_main.append([adj_price, adj_quantity])

    if c_direction == "quote_to_base":
        for p in bids:
            bid_price = float(p[0])
            adj_price = bid_price if bid_price != 0 else 0
            adj_quantity = float(p[1])
            price_list_main.append([adj_price, adj_quantity])

    return price_list_main


def calc_acquired_coin(amount_in, orderbook):
    # depth calculation
    """
        CHALLENGES

        full amount of starting amount can be eaten at the first level (lvl 0)
        some of the amount in can be eaten up by multiple levels
        some coins may not have enough liquidity
    """

    trading_balance = amount_in
    quantity_bought = 0
    acquired_coin = 0
    counts = 0
    amount_bought = 0

    for level in orderbook:
        # extract level price and quantity
        level_price = level[0]
        level_available_quantity = level[1]

        # amount in is <= first level total amount
        if trading_balance <= level_available_quantity:
            quantity_bought = trading_balance
            trading_balance = 0
            amount_bought = quantity_bought * level_price

        # amount in is > than a given level total amount
        if trading_balance > level_available_quantity:
            quantity_bought = level_available_quantity
            trading_balance -= quantity_bought
            amount_bought = quantity_bought * level_price

        # accumulate acquired coin
        acquired_coin += amount_bought

        # exit trade
        if trading_balance == 0:
            return acquired_coin

        # exit if not enough orderbook levels
        counts += 1

        if counts == len(orderbook):
            return 0


def get_depth_from_orderbook(surface_dict, prices):
    depth_1_reformatted_prices = reformatted_orderbook(
        prices["pair_a_ask"], prices["pair_a_bid"], surface_dict["direction_trade_1"])
    depth_2_reformatted_prices = reformatted_orderbook(
        prices["pair_b_ask"], prices["pair_b_bid"], surface_dict["direction_trade_2"])
    depth_3_reformatted_prices = reformatted_orderbook(
        prices["pair_c_ask"], prices["pair_c_bid"], surface_dict["direction_trade_3"])

    # get acquired coins / get depth
    acquired_coin_t1 = calc_acquired_coin(
        surface_dict['starting_amount'], depth_1_reformatted_prices)
    acquired_coin_t2 = calc_acquired_coin(
        acquired_coin_t1, depth_2_reformatted_prices)
    acquired_coin_t3 = calc_acquired_coin(
        acquired_coin_t2, depth_3_reformatted_prices)

    # calculate profit loss (aka real rate)
    profit_loss = acquired_coin_t3 - surface_dict['starting_amount']
    real_rate_perc = (
        profit_loss / surface_dict['starting_amount']) * 100 if profit_loss != 0 else 0

    if real_rate_perc > 0:
        return {
            "profit_loss": profit_loss,
            "real_rate_perc": real_rate_perc,
            "contract_1": surface_dict["contract_1"],
            "contract_2": surface_dict["contract_2"],
            "contract_3": surface_dict["contract_3"],
            "contract_1_direction": surface_dict["direction_trade_1"],
            "contract_2_direction": surface_dict["direction_trade_2"],
            "contract_3_direction": surface_dict["direction_trade_3"],
        }
    else:
        return {}
