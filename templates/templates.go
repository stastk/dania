package templates

import (
	"main/models"
	"strconv"
)

func ShowList(list []models.List) string {
	list_html := `
	<style>
		body{
			padding: 100px 630px;
			margin: 0;
		}
		.list{
			background: #ffffff;
			border-radius: 14px;
			overflow: hidden;
			display: block;
			font-family: Arial, sans-serif;
			border: 1px solid #eee;
			padding: 16px;
		}
		.ingridient{
			display: flex;
			min-height: 24px;
			align-items: center;
		}
		.ingridient:not(:last-of-type){
			margin-bottom: 8px;
		}
		.chk_container{
			padding-right: 4px;
		}
		.name_container{
			padding-left: 4px;
			padding-right: 4px;
		}
		.dots{
			flex: 1;
			padding-left: 4px;
			padding-right: 4px;
		}
		.dots i{
			background-image: linear-gradient(to right, black 33%, rgba(255,255,255,0) 0%);
			background-position: bottom;
			background-size: 3px 1px;
			background-repeat: repeat-x;
			width: 100%;
			height: 1px;
			display: flex;
		}
		.count_container{
			padding-left: 4px;
			padding-right: 4px;	
		}
	</style>
	<div class="list">
		<div class="list_of_ingridients">
	`

	i := 0

	for ingridients_count := range list[0].Ingridients {
		for ingridients_count >= i {
			uid := strconv.Itoa(list[0].Ingridients[i].IngridientVariationId) + strconv.Itoa(list[0].Ingridients[i].IngridientId) + strconv.Itoa(list[0].Ingridients[i].Count)
			list_html += `
			<div class="ingridient">
				<div class="chk_container">
					<label for="v` + uid + `" class="line">
						<input type="checkbox" id="v` + uid + `">
					</label>
				</div>
				<div class="name_container">
					<span class="name">` + list[0].Ingridients[i].VariationName + `</span>
				</div>
				<div class="dots">
					<i></i>
				</div>
				<div class="count_container">
					<strong>` + strconv.Itoa(list[0].Ingridients[i].Count) + `</strong> ` + list[0].Ingridients[i].UnitName + `
				</div>
			</div>
			`
			i++
		}
	}

	list_html += `
	</table>
	`

	if len(list[0].Description) > 0 {
		list_html += `<span class="description">` + list[0].Description + `</span>`
	}

	list_html += `
		</div>
	</body>
	`
	return list_html
}
