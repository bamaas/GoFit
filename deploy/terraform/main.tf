terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }

  backend "azurerm" {
    resource_group_name  = "tfstate"
    storage_account_name = "tfstate12179"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }

  required_version = ">= 1.1.0"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "gofit" {
    name     = "gofit"
    location = "West Europe"
}

resource "azurerm_service_plan" "gofit" {
    name                = "ASP-gofit-ab05"
    location            = azurerm_resource_group.gofit.location
    resource_group_name = azurerm_resource_group.gofit.name
    sku_name = "F1"
    os_type  = "Linux" 
}

resource "azurerm_linux_web_app" "gofit" {
    name                = "gofit-api"
    location            = azurerm_resource_group.gofit.location
    resource_group_name = azurerm_resource_group.gofit.name
    service_plan_id = azurerm_service_plan.gofit.id
    https_only = true

    app_settings = {
        "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "true"
        "PORT" = "80"
        "WEBSITES_PORT" = "8080"
        "USERS" = "[{\"email\": \"demo@gofit.nl\", \"password\": \"gofit123\"}, {\"email\": \"test@gofit.nl\", \"password\": \"gofit123\"}]"
    }

    site_config {
        always_on = false
        cors {
            allowed_origins = ["*"]
        }
    }
}

// Deprecated resource but import doesn't work with the new 'azurerm_static_web_app'.
// To migrate: 'terraform state rm azurerm_static_web_app.gofit && import again with the new resource name'
resource "azurerm_static_site" "gofit" {
    name                = "gofit-frontend"
    resource_group_name = azurerm_resource_group.gofit.name
    location            = azurerm_resource_group.gofit.location
}


// Create function app

resource "azurerm_resource_group" "gofit_functions" {
    name     = "gofit-functions"
    location = "West Europe"
}

resource "azurerm_storage_account" "gofit" {
  name                     = "gofit"
  resource_group_name      = azurerm_resource_group.gofit_functions.name
  location                 = azurerm_resource_group.gofit_functions.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "gofit_functions" {
    name                = "ASP-gofit-functions"
    location            = azurerm_resource_group.gofit_functions.location
    resource_group_name = azurerm_resource_group.gofit_functions.name
    sku_name = "Y1"
    os_type  = "Linux" 
}

resource "azurerm_linux_function_app" "gofit_caffeine" {
    name                = "gofit-caffeine"
    location            = azurerm_resource_group.gofit_functions.location
    resource_group_name = azurerm_resource_group.gofit_functions.name

    service_plan_id = azurerm_service_plan.gofit_functions.id
    storage_account_name       = azurerm_storage_account.gofit.name
    storage_account_access_key = azurerm_storage_account.gofit.primary_access_key

    site_config {}
}