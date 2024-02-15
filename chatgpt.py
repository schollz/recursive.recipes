from openai import OpenAI

import json
import os

client = OpenAI(api_key="")

content = """You are a recipe API. I will give you an recipe and you will give JSON.

The JSON should have the original recipe name, a singular recipe name, and a plural recipe name.

The JSON should have a list of inputs, and a list of outputs. In both cases, the inputs and outputs should contain a list of ingredients. Each ingredient should have a name, a name in singular form, an amount, a type, and a store cost. The amount should have a count and a unit. The type should be one of: natural, plants, animal, grain, spice, synthetic, or fungi. The store cost should be the cost in US dollars of buying the ingredient from a store or restaurant in the quantity needed for the recipe. 

The JSON should also contain a list of directions.

The JSON should also have the total amount of time in hours needed to make the recipe split into two types of times: parallel and serial. When multiplying the time of the recipe, the parallel time does not increase (e.g. the time of baking does not increase if you double the recipe). In contrast, the serial time does increase when multiplying the recipe (the time for cutting vegetables doubles when there are twice as many recipes).

If the recipe I ask for is not a food recipe, you will still try to make a recipe for it. For example, if asked for a recipe for "whole chicken" you will give the ingredients for raising chickens and everything that entails (plots of land, water, feed). For another exampe, if asked for a recipe for "salt" you will give the ingredients "sea water" and "sun". If asked for a recipe for "sun" you will give the ingredients for making a universe. If asked for something like "tomato" then describe how to grow and harvest the tomato. If asked for "sugar cane" describe how to grow and harvest sugar cane. If asked for "water" describe how to combine hydrogen and oxygen in amounts that produces the required amount of water. You should not go further back than the big bang.

Ingredient types should only be one of these: natural (inedible natural products like soil, sun), plants (fruits or vegetables), animal (meat, eggs, milk), grain (flour, rice, legume), spice (salt, pepper, sugar, cinnamon), synthetic (baking soda, baking powder), or fungi (yeast, mushroom, rennet).

For plant recipes, like "potato", describe the processs of using plant seeds to grow the crop and estimate the time for growing the crop and maintaining it and harvesting it. 

For synthetic ingredients, you should describe the process of making the ingredient. For example, if the recipe calls for "baking powder" you should describe the industrical chemical process of making baking powder. If the recipe calls for "baking soda" you should describe the process of mining and refining baking soda. If the recipe calls for "sugar" you should describe the process of refining sugar from sugar cane or sugar beets. If the recipe calls for "salt" you should describe the process of mining and refining salt.

Be specific when asking for milk, and designate whether it is regular milk ("milk") or "cow milk" (which is used when making butter). Also distinguish between live animals and meat, for example use "live chicken" and "chicken" to distinguish between the animal and the meat.

When asked for the sun, describe how a sun is formed inside a galaxy.

When asked for making yeaset, describe how to make it from grapes.

The time (parallel or serial) should only ever be in hours. Even if it is billions of years, you should represent it in hours in scientific notation if needed. 

The costs should always be in dollars. Even if the cost is astronomical you should make a best estimate in dollars, even if it needs scientific notation. The costs should all be the approximate cost of buying the ingredient from the store in the quantity needed for the recipe. The output should be the cost of the product when purchased from a store or restaurant.

The units should be "grams" when it is weight-based and "cups" if it is volume based. You can use other units when nessecary (like "whole" if it an "egg" or "apple pie" that can't be measured in other units).

Ingredients should always be physical entities and cannot be things like "gravity" or "nucleosynthesis". Physical objects, foods, materials, gases, chemicals, are all examples of things that are okay.

Amounts should always be converted to grams or cups if possible. For example, instead of "1 block of cheese", you should say "100.0 grams of cheese". Instead of "1 kg" use "1000.0 grams" Instead of "1 teaspoon" use "0.021 cups". If it is not possible to convert to grams or cups, then use the most appropriate unit (like "1 whole chicken").

Most important, make sure that the directions are clear and that each input required is mentioned at least once in the complete directions.

"""


pizza_json = """{
    "directions": [
        "In a large bowl, combine the flour and salt.",
        "In a separate small bowl, dissolve the yeast in warm water and let it sit for 5 minutes.",
        "Add the yeast mixture and olive oil to the flour mixture.",
        "Stir until a dough forms, then knead the dough on a floured surface for about 5 minutes.",
        "Place the dough in a greased bowl, cover with a damp towel, and let it rise for 1-2 hours or until doubled in size.",
        "Preheat the oven to 475\u00b0F (245\u00b0C).",
        "Punch down the dough, then divide it into desired portions.",
        "On a floured surface, roll out each portion into a round shape.",
        "Place the rolled out dough on a baking sheet or pizza stone.",
        "Spread tomato sauce evenly over the dough, leaving a small border for the crust.",
        "Sprinkle mozzarella cheese generously on top of the sauce.",
        "Add desired toppings, such as pepperoni, mushrooms, or bell peppers.",
        "Bake the pizza for 12-15 minutes or until the crust is golden brown and the cheese is melted and bubbly.",
        "Remove from the oven and let it cool for a few minutes before slicing."
    ],
    "hours": {
        "parallel": 1.5,
        "serial": 2
    },
    "input": [
        {
            "amount": {
                "count": 3,
                "unit": "cups"
            },
            "ingredient": "flour",
            "singular": "flour",
            "store_cost": 0.6,
            "type": "grain"
        },
        {
            "amount": {
                "count": 1,
                "unit": "cup"
            },
            "ingredient": "water",
            "singular": "water",
            "store_cost": 0,
            "type": "natural"
        },
        {
            "amount": {
                "count": 7,
                "unit": "grams"
            },
            "ingredient": "yeast",
            "singular": "yeast",
            "store_cost": 0.35,
            "type": "fungi"
        },
        {
            "amount": {
                "count": 10,
                "unit": "grams"
            },
            "ingredient": "salt",
            "singular": "salt",
            "store_cost": 0.05,
            "type": "spice"
        },
        {
            "amount": {
                "count": 0.125,
                "unit": "cup"
            },
            "ingredient": "olive oil",
            "singular": "olive oil",
            "store_cost": 0.1,
            "type": "natural"
        },
        {
            "amount": {
                "count": 1,
                "unit": "cup"
            },
            "ingredient": "tomato sauce",
            "singular": "tomato sauce",
            "store_cost": 0.5,
            "type": "plant"
        },
        {
            "amount": {
                "count": 0.5,
                "unit": "cup"
            },
            "ingredient": "mozzarella cheese",
            "singular": "mozzarella cheese",
            "store_cost": 1.5,
            "type": "animal"
        }
    ],
    "output": [
        {
            "amount": {
                "count": 1,
                "unit": "whole"
            },
            "ingredient": "pizza",
            "singular": "pizza",
            "store_cost": 9.99,
            "type": "food"
        }
    ],
    "recipe": "pizza",
    "recipe_plural": "pizzas",
    "recipe_singular": "pizza"
}"""

sea_salt_json = """{
    "recipe": "sea salt",
    "recipe_plural": "sea salt",
    "recipe_singular": "sea salt",
    "input": [
        {
            "ingredient": "sea water",
            "singular": "sea water",
            "amount": { "count": 1, "unit": "cups"},
            "type": "natural",
            "store_cost": 0
        },
        {
            "ingredient": "sun",
            "singular": "sun",
            "amount": { "count": 10.0, "unit": "days"},
            "type": "natural",
            "store_cost": 0
        }
    ],
    "output": [
        {
            "ingredient": "sea salt",
            "singular": "salt",
            "amount": { "count": 0.041, "unit": "cups"},
            "type": "natural",
            "store_cost": 0.01
        }    
    ]
    "directions": [
        "Collect seawater from the ocean",
        "Put the seawater in the sunlight and let it evaporate until it becomes salt"
    ],
    "hours": {"serial": 240.0, "parallel": 0.5},
}
"""

chicken_json = """{
    "directions": [
        "Obtain a chicken and provide it with appropriate living conditions, including shelter and space to roam.",
        "Feed the chicken with 100 kilograms of chicken feed and 100 liters of water over a period of approximately 6 months.",
        "Monitor the chicken's health and well-being, ensuring it has a balanced diet and access to fresh water.",
        "Wait for the chicken to lay an egg, which typically takes around 24-26 hours.",
        "Collect the freshly laid egg."
    ],
    "hours": {
        "parallel": 4380.0,
        "serial": 0.2,
    },
    "input": [
        {
            "amount": {
                "count": 1,
                "unit": "whole"
            },
            "ingredient": "chicken",
            "singular": "chicken",
            "store_cost": 6.99,
            "type": "animal"
        },
        {
            "amount": {
                "count": 100.0,
                "unit": "kilograms"
            },
            "ingredient": "chicken feed",
            "singular": "chicken feed",
            "store_cost": 45.0,
            "type": "natural"
        },
        {
            "amount": {
                "count": 100.0,
                "unit": "liters"
            },
            "ingredient": "water",
            "singular": "water",
            "store_cost": 0.001,
            "type": "natural"
        }
    ],
    "output": [
        {
            "amount": {
                "count": 125.0,
                "unit": "whole"
            },
            "ingredient": "egg",
            "singular": "egg",
            "store_cost": 50.00,
            "type": "animal"
        }
    ],
    "recipe": "egg",
    "recipe_plural": "eggs",
    "recipe_singular": "egg"
}"""


def generate_recipe(recipe):
    filename = f"static/data/{recipe}.json"
    if os.path.exists(filename):
        return
    chat_completion = client.chat.completions.create(
        model="gpt-3.5-turbo",
        messages=[
            {"role": "system", "content": content},
            {
                "role": "user",
                "content": "Generate a recipe for 'pizza'",
            },
            {"role": "system", "content": pizza_json},
            {
                "role": "user",
                "content": "Generate a recipe for 'sea salt'",
            },
            {"role": "system", "content": sea_salt_json},
            {
                "role": "user",
                "content": f"Generate a recipe for '{recipe}'",
            },
        ],
    )

    print(chat_completion.choices)
    a = chat_completion.choices[0].message.content
    b = json.loads(a)
    with open(f"static/data/{recipe}.json", "w") as f:
        f.write(json.dumps(b, indent=2, sort_keys=True))
    print(f"Generated {recipe}.")


recipes = json.load(open("static/data/recipes.json", "rb"))

for recipe in recipes:
    try:
        generate_recipe(recipe)
    except Exception as e:
        print(f"Failed to generate {recipe}: {e}")
        continue
